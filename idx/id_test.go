package idx

import (
	"fmt"
	"testing"
)

func TestSnowFlake(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(SnowFlake().Value())
	}
}

func TestID(t *testing.T) {
	t.Run("snowflake", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			fmt.Println(SnowFlake().Value())
		}
	})

	t.Run("sequence", func(t *testing.T) {
		name := "test"
		for i := 0; i < 100; i++ {
			fmt.Println(GetSequence(name).Next())
		}
	})

	t.Run("uuid", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			fmt.Println(UUID())
		}
	})

	t.Run("time", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			fmt.Println(TimeUnix())
		}
	})

	t.Run("time-milli", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			fmt.Println(TimeUnixMilli())
		}
	})

	t.Run("timestamp", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			fmt.Println(Timestamp())
		}
	})
}
