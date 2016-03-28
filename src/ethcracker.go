// Written by @AlexNa

package main

import ( 
    "flag"
    "os"
    "bufio"
    "sync"
    "strings"
    "unicode"
//    "encoding/json"
//    "fmt"
    "time"
    "./crypto"
)

var templates [][]string 

type TEMP_FLAGS struct {
    UseAlways bool
    Capitalize bool
}

var templates_flags []TEMP_FLAGS


// note, that variables are pointers
var pk = flag.String("pk", "", "Private key file")
var t = flag.String("t", "", "Pattern file")
var l = flag.String("l", "", "File with list of variants. If specified, -t is ignored" )
var min_len = flag.Int("min_len", 8, "Minimum password length")
var max_len = flag.Int("max_len", 20, "Maximum password length")
var n_threads = flag.Int("threads", 4, "Number of threads")
var pre_sale = flag.Bool("presale", false, "The key file is the presale JSON")
var keep_order = flag.Bool("keep_order", false, "Keep order of the lines (no permutations)")
var v = flag.Int("v", 1, "Verbosity ( 0, 1, 2 )")
var re = flag.Int("re", 1, "Report every N-th combination")
var start_from = flag.Int("start_from", 0, "Skip first N combinations")
var dump = flag.String("dump", "", "Just output all the possible variants")

var params crypto.CrackerParams
var chans []chan string
var wg sync.WaitGroup
var f_dump *os.File 

func fact( x int) int {
  if x == 0 {
    return 1
  }
  return x * fact( x - 1 )
}

func safe_add( a []string, s string ) []string {
    for _,n := range( a ) {
        if n == s { return a }
    }
    
    return append( a, s )
}


func main() {
    var err error
    
    flag.Parse()
    
    if *dump != "" { 
        *v = 0
        *n_threads = 1
        
        f_dump, err = os.Create( *dump )
        if err != nil { panic( err ) }
        
        defer f_dump.Close() 
    }
    
    
    if *v > 0 {
        println( "------------------------------------------------")
        println( "Ethereum Password Cracker v2.0")
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
    }
    
    if *re <=0 { panic( "wrong -re")}
    
    if *v > 0 { println( "Report every :", *re, "combination") }
    
    params.V = *v
    params.Start_from = *start_from
    params.StartTime = time.Now()
    params.RE = *re
    
    if *pk == "" { panic( "No key file") }
    
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
        err := crypto.LoadKeyFile( &params, *pk, *v)
        if err != nil { panic( err ) }
    }
    
    templates = make( [][]string, 0 )
    templates_flags = make( []TEMP_FLAGS, 0 )

    pl := make( []string, 0 )
    
    if *l != "" {
        f, err := os.Open( *l )
        if err != nil { panic( err ) }

        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            pl = append( pl, scanner.Text() )
        }

        if err := scanner.Err(); err != nil { panic( err ) }
        f.Close()
    } else {
        if *t == "" { panic( "No template file") }
        
        f, err := os.Open( *t )
        if err != nil { panic( err ) }

        scanner := bufio.NewScanner(f)
        for scanner.Scan() {
            
            templ := make( []string, 0 )
            
            tl := strings.Split( scanner.Text(), " " )
            
            for _, n := range( tl ) {
                if n == "" { continue }
                if n == " " { continue }
                
                templ = safe_add( templ, n )                
            }
            
            if len( templ ) >= 0 { 

                for i, _ := range( templ ) {
                    templ[i] = strings.Replace( templ[i], "\\s", " ", -1 )
                }
                
                var tf TEMP_FLAGS
                
                if strings.HasPrefix( templ[0], "~" ) {
                    if len( templ ) == 1 { continue } //nothing but flags...
                    
                    tf.UseAlways = strings.Index( templ[0], "a" ) > 0
                    tf.Capitalize = strings.Index( templ[0], "c" ) > 0
                 
                    templ = templ[1:] //remove the first 
                }
                
                if tf.Capitalize {
                    t := make( []string, 0 )

                    for _, n := range( templ ) {
                        if len( n ) > 0 {
                            t = safe_add( t, string( unicode.ToUpper( []rune(n)[0] ) ) + n[1:] )
                            t = safe_add( t, string( unicode.ToLower( []rune(n)[0] ) ) + n[1:] )
                        }
                        templ = t
                    }                    
                }
                

                templates = append( templates, templ ) 
                templates_flags = append( templates_flags, tf ) 
                
                //println( "templates_flags:", len( templates_flags ) - 1, tf.UseAlways ) 
            }
        }

        if err := scanner.Err(); err != nil { panic( err ) }

        f.Close()
    
        if *v > 0 { println( "Template lines:", len( templates ) ) }
    }
    
    
    //calculate number of variants:
    
    counters := make( []int, len( templates ) + 1 )
    indexes := make( []int, len( templates ) )
    for i := 0; i < len( indexes ); i++ {
        if  templates_flags[i].UseAlways { indexes[i] = 1 }
    }

    if *l != "" {
        params.Total = len( pl )
    } else if *keep_order {
        
        a_exists := false
        
        params.Total = 1;
        for i := 0; i < len( templates ); i++ { 
            if templates_flags[i].UseAlways {
                params.Total = params.Total  * len( templates[i] )  
                a_exists = true
            } else {
                params.Total = params.Total  * ( len( templates[i] ) + 1 ) 
            }
        }
        if !a_exists { params.Total = params.Total - 1 }
    } else {
        if len( templates ) > 20 { panic( "Too many templates. No way you have so much powerful computer...")}
    
        counter: for {
            not_zero := 0
            for _,k := range( indexes ) { if k != 0 { not_zero++ }}

            counters[not_zero]++

            for i := 0; i < len( indexes ); i++ {

                if indexes[i] < len( templates[i] ) {
                    indexes[i] = indexes[i] + 1
                    break;
                } else {
                    if templates_flags[i].UseAlways {
                        indexes[i] = 1
                    } else {
                        indexes[i] = 0
                    }                        
                    if i == len( templates ) - 1 { break counter }
                }
            } 
        }

        for i, c := range( counters ) {
            //println( "counters ", i, c )
            params.Total += c * fact( i )

            if params.Total > 100000000 { panic( "Too many templates. No way you have so much powerful computer...")}

        }
    }
    
    if *v > 0 { println( "Total possible variants:", params.Total) }
    if *v > 0 { println( "---------------- STARTING ----------------------") }

    //main cycle
    if *l != "" {
        for _, np := range( pl ) { 
            if( *n_threads == 1 ) {
                crypto.Test_pass( &params, np, 0 )
            } else {
                chans[ params.N % *n_threads ] <- np
            }        
        }
    } else {
        indexes = make( []int, len( templates ) )
        for i := 0; i < len( indexes ); i++ {
            if  templates_flags[i].UseAlways { indexes[i] = 1 }
        }
        
        main: for {
            
            letters := make( []string, 0 )

            for i := 0; i < len( indexes ); i++ {
                if indexes[i] > 0 { letters = append( letters, templates[i][indexes[i] - 1 ] ) }
            }

            letters_str := ""
            for _, l := range( letters ) { letters_str += "(" + l + ") "}

            if *v > 1 && params.N >= params.Start_from { println( "Selected letters:", letters_str ) }

            if *keep_order {
                test( letters )
            } else {
                AllPermutations( letters, 0 )
            }

            for i := 0; i < len( indexes ); i++ {

                if indexes[i] < len( templates[i] ) {
                    indexes[i] = indexes[i] + 1
                    break;
                } else {
                    if templates_flags[i].UseAlways {
                        indexes[i] = 1
                    } else {
                        indexes[i] = 0
                    }                        
                    if i == len( templates ) - 1 { break main }
                }
            } 
        
        }
    }
    
    //wait for threads to finish
    if *n_threads > 1  {
        for i := 0; i < *n_threads; i++ { 
            chans[i] <- ""
        }
        wg.Wait()        
    }
           
    if *v > 0 { println( ":-( Sorry... password not found") }
} 

func test( l []string ) {
    s := ""
    for _, n := range( l ) { s = s + n }

    if len(s) < *min_len { return }
    if len(s) > *max_len { return }
    
    if *dump != "" {
        f_dump.Write( []byte( s + "\n" ) )
        return
    }    
    
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