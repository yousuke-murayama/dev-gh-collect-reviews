package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/spf13/cobra"
)

var(
 repo string
 author string
 from string
 to string
 label string
)

var rootCmd = &cobra.Command{
 Use:   "collect-reviews",
 // メイン処理
 RunE: func(cmd *cobra.Command, args []string) error {
  repo, err := cmd.Flags().GetString("repo");
  if err != nil {
   return err
  }

  author, err := cmd.Flags().GetString("author");
  if err != nil {
   return err
  }

  from, err := cmd.Flags().GetString("from");
  if err != nil {
   return err
  }

  to, err := cmd.Flags().GetString("to");
  if err != nil {
   return err
  }

  label, err := cmd.Flags().GetString("label");
  if err != nil {
   return err
  }

  repoQuery := fmt.Sprintf("--repo=%s", repo)
  authorQuery := fmt.Sprintf("--author=%s", author)
  dateQuery := fmt.Sprintf("--created=%s..%s", from, to)
  stdOut, _, err := gh.Exec("search", "prs", repoQuery, authorQuery, dateQuery, "--limit=50", "--json", "number")
  if err != nil {
   return err
  }
  var prs []struct { Number int };
  err = json.Unmarshal(stdOut.Bytes(), &prs);
  if err != nil {
   return err
  }

  var comments []struct {Body string; Html_url string;};
  client, err := api.DefaultRESTClient()
  if err != nil {
   return err
  }

  for _, pr := range prs {
   response := []struct {Html_url string; Body string; User struct { Login string }}{}
   err = client.Get(fmt.Sprintf("repos/%s/pulls/%d/comments", repo, pr.Number), &response)

   if err != nil {
    return err
   }
   // FIXME: 条件分岐が汚いのでリファクタしたい
   if len(response) > 0 {
    for _, res := range response {
     if res.User.Login != author {
      if(label != "") {
       if(strings.Contains(res.Body, label)) {
        comments = append(comments, struct{ Body string; Html_url string; }{Body: res.Body, Html_url: res.Html_url})
       }
      } else {
       comments = append(comments, struct{ Body string; Html_url string; }{Body: res.Body, Html_url: res.Html_url})
      }
     }
    }
   }
  }

  for i, comment := range comments {
   if(i == 0) {
    fmt.Println("--------------------------------------------------------------------------------------")
   }
   fmt.Println("No." + fmt.Sprintf("%d", i+1))
   fmt.Println(comment.Body)
   fmt.Println(comment.Html_url)
   fmt.Println("--------------------------------------------------------------------------------------")
  }
  fmt.Println("total: " + fmt.Sprintf("%d", len(comments)));
  return nil
 },
}

func main() {
 err := rootCmd.Execute()
 if err != nil {
  os.Exit(1)
 }
}

func init() {
 rootCmd.Flags().StringVarP(&repo, "repo", "r", "", "リポジトリ名")
 rootCmd.Flags().StringVarP(&author, "author", "a", "", "PRの作成者")
 rootCmd.Flags().StringVarP(&from, "from", "f", "", "PRの作成日")
 rootCmd.Flags().StringVarP(&to, "to", "t", "", "PRの作成日")
 rootCmd.Flags().StringVarP(&label, "label", "l", "", "レビューコメントのラベル")

 rootCmd.MarkFlagRequired("repo")
 rootCmd.MarkFlagRequired("author")
 rootCmd.MarkFlagRequired("from")
 rootCmd.MarkFlagRequired("to")

 rootCmd.Flags().SortFlags = false
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
