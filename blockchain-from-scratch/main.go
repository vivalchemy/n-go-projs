package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"fmt"
	"github.com/gorilla/mux"
)

type Block struct {
	Pos       int
	Timestamp time.Time
	Hash      string
	PrevHash  string
	Data      BookCheckout
}

type BookCheckout struct {
	IsGenesis    bool   `json:"is_genesis"`
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type BlockChain struct {
	blocks []*Block
}

var blockChain = new(BlockChain)

func (b *Block) generateHash() string {
	h := sha256.New()
	io.WriteString(h, fmt.Sprint(b.Pos)+b.Timestamp.String()+b.PrevHash+b.Data.BookID+b.Data.User+b.Data.CheckoutDate)
	return hex.EncodeToString(h.Sum(nil))
}

func (b *Block) validateHash() bool {
	return b.Hash == b.generateHash()
}

func CreateBlock(prevBlock *Block, checkoutItem *BookCheckout) *Block {
	block := &Block{}

	block.Pos = prevBlock.Pos + 1
	block.PrevHash = prevBlock.Hash
	block.Timestamp = time.Now()
	block.Data = *checkoutItem

	block.Hash = block.generateHash()

	return block
}

func validBlock(block *Block, prevBlock *Block) bool {
	// there can be only one genesis block which is created when the program starts
	// so the user can't create a new genesis block
	if block.Data.IsGenesis {
		return false
	}
	if block.Pos != prevBlock.Pos+1 {
		return false
	}
	if block.PrevHash != prevBlock.Hash {
		return false
	}
	if block.Timestamp.Before(prevBlock.Timestamp) {
		return false
	}
	if !block.validateHash() {
		return false
	}

	return true
}

func (b *BlockChain) AddBlock(checkoutItem *BookCheckout) {
	// create new block
	prevBlock := b.blocks[len(b.blocks)-1]

	block := CreateBlock(prevBlock, checkoutItem)

	if validBlock(block, prevBlock) {
		b.blocks = append(b.blocks, block)
	}

}

func createNewBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON for book required fields are: title, author, publish_date, isbn"))
	}

	if book.Title == "" || book.Author == "" || book.PublishDate == "" || book.ISBN == "" {
		log.Println("Invalid book")

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid book required fields are: title, author, publish_date, isbn"))
		return
	}

	h := md5.New()
	io.WriteString(h, book.ISBN+book.PublishDate) // data to hash
	book.ID = hex.EncodeToString(h.Sum(nil))

	// alternative to MarshalIndent but the content-type is not working correctly
	// w.WriteHeader(http.StatusCreated)
	// w.Header().Set("Content-Type", "application/json")
	//
	// if err := json.NewEncoder(w).Encode(book); err != nil {
	// 	log.Println(err)
	//
	// 	w.Header().Set("Content-Type", "text/plain")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Error marshalling book"))
	// 	return
	// }

	resp, err := json.MarshalIndent(book, "", "  ")
	// will work the same as the api clients will format the json using content-type
	// resp, err := json.Marshal(book, "", "  ")

	if err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling book"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(resp)
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON for checkout item(block not created) required fields are: book_id, user, checkout_date"))
	}

	blockChain.AddBlock(&checkoutItem)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	fmt.Println(len(blockChain.blocks))

	if err := json.NewEncoder(w).Encode(checkoutItem); err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling checkout item"))
		return
	}
}

func getBlockChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(blockChain.blocks); err != nil {
		log.Println(err)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error marshalling block chain"))
		return
	}
}

func GenesisBlock() *Block {
	return CreateBlock(&Block{}, &BookCheckout{IsGenesis: true})
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{
			GenesisBlock(),
		},
	}
}

func main() {
	// create a new blockchain
	blockChain = NewBlockChain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockChain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", createNewBook).Methods("POST")

	go func() {
		for {
			fmt.Println("----------------------------------------------------------")
			for _, block := range blockChain.blocks {
				fmt.Println("PrevHash: ", block.PrevHash)
				bytes, _ := json.MarshalIndent(block, "", "  ")
				fmt.Println(string(bytes))
				fmt.Println("BlockHash: ", block.Hash)
				fmt.Println("------------")
			}
			time.Sleep(time.Second * 10)
		}
	}()

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
