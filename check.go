package main

func Checklocation(){
	for i:=0 ; i<10 ; i++ {
		for j:=0 ; j<10 ; j++ {
			if Field[i][j].busy==true{
				if i!=0&&j!=0{
					Field[i-1][j-1].access=false
				}else if i!=0&&j!=9{
					Field[i-1][j+1].access=false
				}else if j!=0&&i!=9{
					Field[i+1][j-1].access=false
				}else if j!=9&&i!=9{
					Field[i+1][j+1].access=false
				}
			}
			if Field[i+1][j].busy==true&&i!=9{
				if Field[i-1][j].busy==false&&i!=0{
					Field[i-1][j].access=false
				}
				if j!=0{
					Field[i][j-1].access=false
				}
				if j!=9{
					Field[i][j+1].access=false
				}
				
			}
			if Field[i-1][j].busy==true&&i!=0{
				if Field[i+1][j].busy==false&&i!=9{
					Field[i+1][j].access=false
				}
				if j!=0{
					Field[i][j-1].access=false
				}
				if j!=9{
					Field[i][j+1].access=false
				}
			}
			if Field[i][j+1].busy==true&&j!=9{
				if i!=0{
					Field[i-1][j].access=false
				}
				if i!=9{
					Field[i+1][j].access=false
				}
				if Field[i][j-1].busy==false&&j!=0{
					Field[i][j+1].access=false
				}
			}
			if Field[i][j-1].busy==true&&j!=0{
				if i!=0{
					Field[i-1][j].access=false
				}
				if i!=9{
					Field[i+1][j].access=false
				}
				if Field[i][j+1].busy==false&&j!=9{
					Field[i][j+1].access=false
				}
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