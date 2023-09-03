/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/meowalien/RabbitGather-golang.git/internal/lib/config"
	"github.com/meowalien/RabbitGather-golang.git/internal/lib/graceful_shutdown"
	"github.com/meowalien/RabbitGather-golang.git/internal/lib/server/grpc"
	"github.com/meowalien/RabbitGather-golang.git/internal/module/webinterest"
)

// webInterestCmd represents the webInterest command
var webInterestCmd = &cobra.Command{
	Use:   "web-interest",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gracefulShutdown := graceful_shutdown.NewGracefulShutdown(graceful_shutdown.ShutdownDeadLine(time.Second * 10))

		conf := config.ParseConfig()
		grpcServer := grpc.New(conf.GRPCServer)

		webinterest.New(grpcServer)

		err := grpc.ListenAndServe(gracefulShutdown, grpcServer, conf.GRPCServer.Port)
		if err != nil {
			panic(err)
		}
		fmt.Println("webInterest called")
	},
}

func init() {
	rootCmd.AddCommand(webInterestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webInterestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webInterestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
