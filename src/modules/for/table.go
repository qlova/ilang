package f

import "github.com/qlova/ilang/src"

func ScanForTable(ic *ilang.Compiler, table, strindex, v string) {
        var array_table = ic.Tmp("array_table")
        var bi = ic.Tmp("bucket_index")
        var condition = ic.Tmp("condition")
        var i = ic.Tmp("index")
       // var hash = ic.Tmp("hash")
        
        var vo = v
        if ic.ExpressionType.SubType.Push == "SHARE" {
            v += "_address"
        } 
        
        ic.Assembly(`
IF 1
PUSH `+table+`
HEAP
GRAB `+array_table+`
VAR `+i+`
ADD `+i+` 0 1
VAR `+bi+`
ADD `+bi+` 0 1
VAR `+condition+`
LOOP
    #Find next table value.
    LOOP
        PLACE `+array_table+`
        PUSH `+i+`
        GET i_pointer
        
        IF i_pointer
            PUSH i_pointer
            HEAP
            GRAB i_bucket
        
            PLACE i_bucket
            PUSH `+bi+`
            GET i_value
            PUSH i_value
            
            SUB `+condition+` `+bi+` 1
            PUSH `+condition+`
            GET i_index
            PUSH i_index
            PUSH 255
            
            ADD `+bi+` `+bi+` 2
            SGE `+condition+` `+bi+` #i_bucket
            IF `+condition+`
                ADD `+bi+` 0 1
                ADD `+i+` `+i+` 1
            END
            BREAK
        ELSE
            ADD `+bi+` 0 1
            ADD `+i+` `+i+` 1
            SGE `+condition+` `+i+` #`+array_table+`
            IF `+condition+`
                BREAK
            END
        END
        
    REPEAT
    SGE `+condition+` `+i+` #`+array_table+`
    IF `+condition+`
        BREAK
    END
    RUN i_unhash
    GRAB `+strindex+`
    PULL `+v+`
    
    
    
    
        `)
        
        
    if ic.ExpressionType.SubType.Push == "SHARE" {
        ic.Assembly("PUSH ", v)
        ic.Assembly("HEAP")
        ic.Assembly("GRAB ", vo)
    }
    
    ic.GainScope()
    ic.GainScope()
    ic.SetFlag(ForLoop)
    
    ic.SetVariable(strindex, ilang.Text)
    ic.SetVariable(vo, *ic.ExpressionType.SubType)
    ic.SetVariable(vo+".", ilang.Protected)
    
}
