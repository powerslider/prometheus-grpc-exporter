package prometheus

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	//nolint:staticcheck
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"
)

// DecodeWriteRequest from an io.Reader into a prompb.WriteRequest, handling
// snappy decompression.
func decodeWriteRequest(r io.Reader) (*prompb.WriteRequest, error) {
	compressed, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	reqBuf, err := snappy.Decode(nil, compressed)
	if err != nil {
		return nil, err
	}

	var req prompb.WriteRequest
	if err := proto.Unmarshal(reqBuf, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func RemoteWriteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := decodeWriteRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var metrics []prompb.TimeSeries
	for _, ts := range req.Timeseries {
		var samples []prompb.Sample
		for _, s := range ts.Samples {
			if math.IsNaN(s.Value) {
				continue
			}
			samples = append(samples, s)
		}
		metrics = append(metrics, prompb.TimeSeries{Labels: ts.Labels, Samples: samples})
	}

	resp, err := json.Marshal(metrics)
	//log.Println(string(resp))
	if err != nil {
		log.Fatal("error serializing metrics response: ", err)
		return
	}

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("metrics server error: ", err)
	}

	//// Listen for incoming connections.
	//c, err := net.Dial("tcp", "goreplay:8070")
	//if err != nil {
	//	log.Println("Error listening:", err.Error())
	//	os.Exit(1)
	//}
	//// Close the listener when the application closes.
	//defer c.Close()
	//log.Println("dialed TCP client on ", "goreplay:8070")
	//for {
	//	//read in input from stdin
	//	respReader := bytes.NewReader(resp)
	//	reader := bufio.NewReader(respReader)
	//	fmt.Print("Text to send: ")
	//	text, _ := reader.ReadString('\n')
	//	log.Println(text)
	//
	//	//send to socket
	//	fmt.Fprint(c, text+"\n")
	//	break
	//}
}
