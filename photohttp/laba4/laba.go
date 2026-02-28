package main

import (
	"fmt"
	"time"
)

type Channel struct {
	name        string
	capacity    int
	currentLoad int
	overloads   int
}

type NetworkStats struct {
	totalTraffic   int
	overloadEvents int
}

func createChannels() []Channel {
	channels := []Channel{
		{"CH-1", 100, 0, 0},
		{"CH-2", 120, 0, 0},
		{"CH-3", 80, 0, 0},
	}
	return channels
}

func generateTraffic() []int {
	// список "пакетів" трафіку
	traffic := []int{80, 50, 120, 40, 60, 110}
	return traffic
}

// findLeastLoadedChannel повинен знаходити канал з найменшим currentLoad
func findLeastLoadedChannel(channels []Channel) int {
	minIndex := 0
	minLoad := (channels)[0].currentLoad

	for i := 1; i < len(channels); i++ {
		if (channels)[i].currentLoad < minLoad {
			minLoad = (channels)[i].currentLoad
			minIndex = i
		}
	}

	return minIndex
}

// distributeTraffic розподіляє трафік по каналах
func distributeTraffic(channels []Channel, traffic []int) NetworkStats {
	stats := NetworkStats{}

	fmt.Println("Starting traffic distribution...")
	fmt.Println("Number of channels:", len(channels))
	fmt.Println("Traffic batches:", len(traffic))

	// симулюємо кілька тактів
	for i := 0; i < len(traffic); i++ {
		amount := traffic[i]
		fmt.Println("\nTick", i+1, "incoming traffic:", amount)

		idx := findLeastLoadedChannel(channels)

		ch := &channels[idx]

		// додаємо трафік
		ch.currentLoad += amount

		// рахуємо загальний трафік
		stats.totalTraffic += amount

		// перевіряємо перевантаження
		if ch.currentLoad > ch.capacity {
			fmt.Println("WARNING: overload on channel", ch.name)
			stats.overloadEvents++
			ch.overloads++
		}

		fmt.Println("Channel", ch.name, "load:", ch.currentLoad, "/", ch.capacity)
	}

	return stats
}

// findMostOverloadedChannel має знаходити канал з найбільшою кількістю overloads
func findMostOverloadedChannel(channels []Channel) *Channel {
	most := &channels[0]

	for i := 1; i < len(channels); i++ {
		if channels[i].overloads > most.overloads {
			most = &channels[i]
		}
	}

	return most
}

func printChannelState(channels *[]Channel) {
	fmt.Println("\n=== Channels state ===")
	for i := 0; i < len(*channels); i++ {
		ch := (*channels)[i]
		fmt.Println(ch.name, "- load:", ch.currentLoad, "/", ch.capacity, "overloads:", ch.overloads)
	}
}

func printSummary(stats NetworkStats, mostOverloaded *Channel) {
	fmt.Println("\n=== Network summary ===")
	fmt.Println("Total traffic:", stats.totalTraffic)
	fmt.Println("Overload events:", stats.overloadEvents)
	fmt.Println("Most overloaded channel:", mostOverloaded.name, "with", mostOverloaded.overloads, "overloads")

	if stats.overloadEvents == 0 {
		fmt.Println("Status: network is stable.")
	} else {
		fmt.Println("Status: network is overloaded, optimization required.")
	}
}

func main() {
	fmt.Println("Network traffic simulation started at", time.Now())

	channels := createChannels()
	traffic := generateTraffic()

	stats := NetworkStats{}

	stats = distributeTraffic(channels, traffic)

	printChannelState(&channels)

	mostOverloaded := findMostOverloadedChannel(channels)

	printSummary(stats, mostOverloaded)
}
