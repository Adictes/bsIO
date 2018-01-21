package main

func Checklocation(){
	for i:=0 ; i<10 ; i++ {
		for j:=0 ; j<10 ; j++ {
			if Field[i][j].busy==true{
				Field[i+1][j+1].access=false
				Field[i+1][j-1].access=false
				Field[i-1][j+1].access=false
				Field[i-1][j-1].access=false
			}
			if Field[i+1][j].busy==true{
				Field[i-1][j].access=false
				Field[i][j-1].access=false
				Field[i][j+1].access=false
				
			}
			if Field[i-1][j].busy==true{
				Field[i+1][j].access=false
				Field[i][j+1].access=false
				Field[i][j-1].access=false
			}
			if Field[i][j+1].busy==true{
				Field[i+1][j].access=false
				Field[i-1][j].access=false
				Field[i][j].access=false
			}
			if Field[i][j-1].busy==true{
				Field[i+1][j].access=false
				Field[i-1][j].access=false
				Field[i][j+1].access=false
			}
		}
	}
	for i:=0 ; i<10 ; i++ {
		for j:=0 ; j<10 ; j++ {
			if (Field[i][j].busy==true)&&(Field[i][j].access==false){
				return
			}	
		}
	}
}