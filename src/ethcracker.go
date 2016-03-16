package main

import ( 
    "flag"
    "os"
    "bufio"
//    "fmt"
    "strings"
//    "encoding/json"

    "./crypto"
)

var templates [][]string 


// note, that variables are pointers
var pk = flag.String("pk", "", "Private key file")
var t = flag.String("t", "", "Pattern file")
var min_len = flag.Int("min_len", 8, "Minimum password length")
var max_len = flag.Int("max_len", 20, "Maximum password length")
var n_threads = flag.Int("threads", 4, "Number of threads")

var params crypto.CrackerParams
var chans []chan string


func main() {
    flag.Parse()
    
    println( "------------------------------------------------")
    println( "Ethereum Password Cracker")
    println( "Author: @AlexNa ")
    println( "------------------------------------------------")
    println( "Private Key File:", *pk )
    println( "Template File:", *t )
    println( "Minimum password length", *min_len )
    println( "Maximum password length", *max_len )
    println( "Number of threads", *n_threads )
    
    if *n_threads < 1 || *n_threads > 32 { panic( "Wroong muber of threads ")}
    
    chans = make( []chan string, *n_threads )
    for i := 0; i < *n_threads; i++ { 
        chans[i] = make( chan string ) 
        
        go func( index int ) {
            
            for {
                s := <- chans[ index ]
                crypto.Test_pass( &params, s, index )
            }
            
        } ( i )
    }
    
    
    err := crypto.LoadKeyFile( &params, *pk)
    if err != nil { panic( err ) }
    
    templates = make( [][]string, 0 )
    f, err := os.Open( *t )
    if err != nil { panic( err ) }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        tl := strings.Split( scanner.Text(), " " )
        if len( tl ) >= 0 { templates = append( templates, tl ) }
    }

    if err := scanner.Err(); err != nil { panic( err ) }

    f.Close()
    
    println( "Template lines:", len( templates ) )
    
    //main cycle
    indexes := make( []int, len( templates ) )

    v := 0
    main: for {
        for i := 0; i < len( indexes ); i++ {
            
            if indexes[i] < len( templates[i] ) {
                indexes[i] = indexes[i] + 1
                break;
            } else {
                indexes[i] = 0
                if i == len( templates ) - 1 { break main }
            }
        } 
        
        v++
        
        letters := make( []string, 0 )
        
        for i := 0; i < len( indexes ); i++ {
            if indexes[i] > 0 { letters = append( letters, templates[i][indexes[i] - 1 ] ) }
        }
        
        AllPermutations( letters, 0 )
    }
           
    println( ":-( Sorry... password not found")           
} 

func test( l []string ) {
    s := ""
    for _, n := range( l ) { s = s + n }

    if len(s) < *min_len { return }
    if len(s) > *max_len { return }
    
    
    chans[ params.N % *n_threads ] <- s
}

func AllPermutations( l []string, index int ) {
    
    if index >= len( l ) - 1 { 
        test( l ) 
        return;
    }
    
    r := len( l ) - index
    
    AllPermutations( l, index + 1 )
    for j := 1; j < r; j++ {
        //swap i and i + j
        tmp := l[index]
        l[index] = l[index + j]
        l[index + j] = tmp
        AllPermutations( l, index + 1 )
    }
}