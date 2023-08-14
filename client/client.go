package main

import (
	"bufio"
	"context"
	"fmt"
	proto "github.com/robinbraemer/cacher/cacher"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// readCommands reads commands from the user and executes them.
// It takes a CacheClient and a Reader as inputs and does not return anything.
func readCommands(c proto.CacheClient, rd *bufio.Reader) {
	ctx := context.Background()

cmds:
	for {
		fmt.Print("-> ")
		in, _ := rd.ReadString('\n')
		// convert CRLF to LF
		in = strings.ReplaceAll(in, "\n", "")

		args := strings.Split(in, " ")

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "":
		case "set":
			if len(args) < 3 {
				fmt.Println("set <key> <string>")
				continue
			}
			key := args[1]
			val := []byte(strings.Join(args[2:], " "))
			if _, err := c.Set(ctx, &proto.SetRequest{Entry: &proto.Entry{Key: key, Val: val}}); err != nil {
				fmt.Println(err)
			}
			fmt.Printf("Key '%s' now '%s'\n", key, val)
		case "get":
			if len(args) < 2 {
				fmt.Println("get <key>")
				continue
			}
			key := args[1]
			rep, err := c.Get(ctx, &proto.GetRequest{Key: key})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Key '%s' is '%s'\n", key, string(rep.Val))
		case "del":
			if len(args) < 2 {
				fmt.Println("del <key>")
				continue
			}
			key := args[1]
			_, err := c.Del(ctx, &proto.DelRequest{Key: key})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("Key '%s' deleted\n", key)
		case "all":
			stream, err := c.All(ctx, &proto.AllRequest{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			var counter int
			for {
				rep, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					fmt.Println(err)
					continue cmds
				}
				counter++
				fmt.Printf("Key '%s' is '%s'\n", rep.Entry.Key, rep.Entry.Val)
			}
			fmt.Println("Got", counter, "entries")
		case "fill-slow":
			if len(args) < 2 {
				fmt.Println("slow-fill <count>")
				continue
			}
			count, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			t := time.Now()
			fmt.Printf("Filling %d entities without reusing the connection...\n", count)
			fillWithoutConnReuse(ctx, count)
			fmt.Printf("Done! Took %s\n", time.Since(t))
		case "fill":
			if len(args) < 2 {
				fmt.Println("fill <count>")
				continue
			}
			count, err := strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				continue
			}
			t := time.Now()
			fmt.Printf("Filling %d entities with reusing the connection...\n", count)
			fillWithConnReuse(ctx, c, count)
			fmt.Printf("Done! Took %s\n", time.Since(t))
		case "clear":
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		case "empty":
			stream, err := c.All(ctx, &proto.AllRequest{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			for {
				resp, err := stream.Recv()
				if err == io.EOF {
					continue cmds
				}
				_, err = c.Del(ctx, &proto.DelRequest{Key: resp.Entry.Key})
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Deleted '%s'\n", resp.Entry.Key)
			}
		default:
			fmt.Println("Commands:")
			fmt.Println("set <key> <string> - Sets an entry")
			fmt.Println("get <key> - Gets an entry")
			fmt.Println("del <key> - Deletes an entry")
			fmt.Println("all - Get all entries")
			fmt.Println("empty - Empty the cache")
			fmt.Println("clear - Clear console")
			fmt.Println("fill <count> - Fills the cache and benefits from HTTP/2 connection reuse")
			fmt.Println("slow-fill <count> - Fills the cache without connection reuse")
		}
	}
}

// main starts the client and reads commands from the user.
// It does not take any inputs or return anything.
func main() {
	cc, err := grpc.Dial(":50001", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()
	readCommands(proto.NewCacheClient(cc), bufio.NewReader(os.Stdin))
}
