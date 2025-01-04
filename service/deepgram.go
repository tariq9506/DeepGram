package websocket

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/api/listen/v1/websocket/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	interfacesv1 "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces/v1"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"
)

type MyCallback struct{}

func (cb MyCallback) Open(or *interfaces.OpenResponse) error {
	fmt.Println("Connection opened:", or.Type)
	return nil
}

func (cb MyCallback) Message(mr *interfaces.MessageResponse) error {
	if len(mr.Channel.Alternatives) > 0 {
		transcript := mr.Channel.Alternatives[0].Transcript
		if transcript != "" {
			fmt.Println("Transcript:", transcript)
		}
	}
	return nil
}

func (cb MyCallback) Metadata(md *interfaces.MetadataResponse) error {
	fmt.Println("Metadata received:", md)
	return nil
}

func (cb MyCallback) SpeechStarted(ssr *interfaces.SpeechStartedResponse) error {
	fmt.Println("Speech started:", ssr)
	return nil
}

func (cb MyCallback) UtteranceEnd(ur *interfaces.UtteranceEndResponse) error {
	fmt.Println("Utterance ended:", ur)
	return nil
}

func (cb MyCallback) Close(cr *interfaces.CloseResponse) error {
	fmt.Println("Connection closed:", cr)
	return nil
}

func (cb MyCallback) Error(er *interfaces.ErrorResponse) error {
	fmt.Println("Error occurred:", er.ErrMsg)
	return nil
}

func (cb MyCallback) UnhandledEvent(byData []byte) error {
	fmt.Println("Unhandled event:", string(byData))
	return nil
}

func ConnectTODeepGram() {
	ctx := context.Background()
	apiKey := "YOUR_DEEPGRAM_API_KEY"

	// // Set up client options
	clientOptions := &interfacesv1.ClientOptions{
		EnableKeepAlive: true,
	}

	// Configure transcription options
	transcriptOptions := &interfacesv2.LiveTranscriptionOptions{
		Language:    "en-US",
		Punctuate:   true,
		Model:       "nova-2",
		SmartFormat: true,
	}

	// Implement the callback
	callback := MyCallback{}

	// Create a new Deepgram LiveTranscription client
	dgClient, err := client.NewWebSocket(ctx, apiKey, clientOptions, transcriptOptions, callback)
	if err != nil {
		log.Println("ERROR creating LiveTranscription connection:", err)
		return
	}

	// Connect to the WebSocket
	if !dgClient.Connect() {
		log.Println("Failed to connect to Deepgram WebSocket")
		os.Exit(1)
	}

	// Simulate a stream for a duration (e.g., 30 seconds)
	go func() {
		time.Sleep(30 * time.Second)
		dgClient.Stop()
	}()

	// Block main thread until connection is closed
	select {}
}
