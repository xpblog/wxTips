package main

func main() {
	go webStarter()
	go wxBotStarter()
	select {}
}
