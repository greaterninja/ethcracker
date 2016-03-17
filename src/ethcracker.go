package main

import ( 
    "flag"
    "os"
    "bufio"
    "sync"
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
var pre_sale = flag.Bool("presale", false, "The key file is the presale JSON")
var keep_order = flag.Bool("keep_order", false, "Keep order of the lines (no permutations)")
var v = flag.Int("v", 0, "Verbosity ( 0, 1, 2 )")
var start_from = flag.Int("start_from", 0, "Skip first N combinations")

var params crypto.CrackerParams
var chans []chan string
var wg sync.WaitGroup

func fact( x int) int {
  if x == 0 {
    return 1
  }
  return x * fact( x - 1 )
}

func main() {
    flag.Parse()
    
    println( "------------------------------------------------")
    println( "Ethereum Password Cracker v1.2")
    println( "Author: @AlexNa ")
    println( "------------------------------------------------")
    println( "Private Key File:", *pk )
    println( "Template File:", *t )
    println( "Verbosity:", *v )
    println( "Minimum password length:", *min_len )
    println( "Maximum password length:", *max_len )
    println( "Number of threads:", *n_threads )
    println( "Presale file:", *pre_sale )
    println( "Keep order:", *keep_order )
    
    params.V = *v
    params.Start_from = *start_from
    
    if *pk == "" { panic( "No key file") }
    if *t == "" { panic( "No template file") }
    
    if *n_threads < 1 || *n_threads > 32 { panic( "Wrong muber of threads ")}
    
    if *n_threads > 1 {
        wg.Add( *n_threads )
        chans = make( []chan string, *n_threads )
        for i := 0; i < *n_threads; i++ { 
            chans[i] = make( chan string ) 

            go func( index int ) {

                for {
                    s := <- chans[ index ]
                    
                    if s == "" { wg.Done(); break; }
                    
                    crypto.Test_pass( &params, s, index )
                }

            } ( i )
        }
    }
    
    if *pre_sale {
        err := crypto.LoadPresaleFile( &params, *pk)
        if err != nil { panic( err ) }
    } else {
        err := crypto.LoadKeyFile( &params, *pk)
        if err != nil { panic( err ) }
    }
    
    templates = make( [][]string, 0 )
    f, err := os.Open( *t )
    if err != nil { panic( err ) }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        tl := strings.Split( scanner.Text(), " " )
        if len( tl ) >= 0 { 
            
            for i, _ := range( tl ) {
                tl[i] = strings.Replace( tl[i], "\\s", " ", -1 )
            }
            
            templates = append( templates, tl ) 
        }
    }

    if err := scanner.Err(); err != nil { panic( err ) }

    f.Close()
    
    println( "Template lines:", len( templates ) )
    
    
    if len( templates ) > 20 { panic( "Too many templates. No way you have so much powerful computer...")}
    
    //calculate number of variants:
    params.Total = 0
    
    if *keep_order {
        params.Total = 1;
        for _, l := range( templates ) { 
            params.Total *= len( l ) + 1
        }
        params.Total = params.Total - 1
        
    } else {
        n := len( templates )

        for k := 1; k <= n; k++ {
            params.Total += fact( n ) / fact( n - k )
        }

        n1 := 0;
        for k := 1; k <= n - 1; k++ {
            n1 += fact( n - 1 ) / fact( n - 1 - k )
        }

        n1 = params.Total - n1 
        for _, l := range( templates ) { 
            params.Total += params.Total * n1 * ( len( l ) - 1 ) / params.Total
        }
    }
    
    println( "Total possible variants:", params.Total)
    
    
    
    println( "---------------- STARTING ----------------------")

    //main cycle
    indexes := make( []int, len( templates ) )
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
        
        letters := make( []string, 0 )
        
        for i := 0; i < len( indexes ); i++ {
            if indexes[i] > 0 { letters = append( letters, templates[i][indexes[i] - 1 ] ) }
        }
        
        letters_str := ""
        for _, l := range( letters ) { letters_str += "(" + l + ") "}
        
        if *v > 0 && params.N >= params.Start_from { println( "Selected letters:", letters_str ) }
        
        if *keep_order {
            test( letters )
        } else {
            AllPermutations( letters, 0 )
        }
    }
    
    //wait for threads to finish
    if *n_threads > 1  {
        for i := 0; i < *n_threads; i++ { 
            chans[i] <- ""
        }
        wg.Wait()        
    }

           
    println( ":-( Sorry... password not found")           
} 

func test( l []string ) {
    s := ""
    for _, n := range( l ) { s = s + n }

    if len(s) < *min_len { return }
    if len(s) > *max_len { return }
    
    
    if( *n_threads == 1 ) {
        crypto.Test_pass( &params, s, 0 )
    } else {
        chans[ params.N % *n_threads ] <- s
    }
}


func makecopy( l []string ) []string {
    nl := make( []string, len( l ) )
    copy( nl, l )
    return nl
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
        AllPermutations( makecopy( l ), index + 1 )
    }
}