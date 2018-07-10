

//this.value_arr=arr;        //不要包含重复值
var index=-1;          //当前选中的DIV的索引
var search_value="";   //保存当前搜索的字符


window.baidu=new Array(); 
function searchSuggest() {
    if (event.keyCode != 40 && event.keyCode != 38 && event.keyCode != 13) {
        var head = document.getElementsByTagName("head")[0];
        var str = encodeURIComponent(document.getElementById('q').value);
        search_value = str;

        var url = 'http://suggestion.baidu.com/su?wd=' + str + '&t=' + Math.round(new Date().getTime()/1000);
        load_script(url, function () {
            window.baidu.sug = function (params) {
                var list = params.s;
                var ss = document.getElementById('suggest')
                ss.innerHTML = '';
                for (i = 0; i < list.length - 1; i++) {
                    var suggest = '<div onmouseover="javascript:suggestOver(this);" ';
                    suggest += 'onmouseout="javascript:suggestOut(this);" ';
                    suggest += 'onclick="javascript:setSearch(this.innerHTML);" ';
                    suggest += 'class="suggest_link">' + list[i] + '</div>';
                    ss.innerHTML += suggest;
                }
            };
        });

    }
}
function load_script(url, callback){ 
    var head = document.getElementsByTagName('head')[0]; 
    var script = document.createElement('script'); 
    script.type = 'text/javascript'; 
    script.src = url; 
    script.onload = script.onreadystatechange = function(){ 
        if((!this.readyState || this.readyState === "loaded" || this.readyState === "complete")){ 
            callback && callback(); 
            script.onload = script.onreadystatechange = null; 
            if ( head && script.parentNode ) { 
                head.removeChild( script ); 
            } 
        } 
    }; 
    head.insertBefore( script, head.firstChild ); 
} 
 
//Mouse over function 
function suggestOver(div_value) { 
div_value.className = 'suggest_link_over'; 
} 
//Mouse out function 
function suggestOut(div_value) { 
div_value.className = 'suggest_link'; 
} 
//Click function 
function setSearch(value) { 
document.getElementById('q').value = value; 
document.getElementById('suggest').innerHTML = ''; 
} 



//键盘上下按键的响应事件
   

 
document.onkeydown=function(event){
    var e = event || window.event || arguments.callee.caller.arguments[0];
    var length = $('#suggest').children().length;;
    var obj=document.getElementById('q');
        //光标键"↓"
        if(event.keyCode==40){
            ++index;
            if(index>length){
                index=0;
            }else if(index==length){
                obj.value=search_value;
            }
            changeClassname(length);
        }
        //光标键"↑"
        else if(event.keyCode==38){
            index--;
            if(index<-1){
                index=length - 1;
            }else if(index==-1){
                obj.value=search_value;
            }
            changeClassname(length);
        }
        //回车键
        else if(event.keyCode==13){
            autoObj.innerHTML = '';
            index=-1;
        }else{
            index=-1;
        }
}

   //更改classname
function changeClassname (length){
    var obj=document.getElementById('q');
    var autoObj=document.getElementById('suggest');//DIV的根节点
        for(var i=0;i<length;i++){
            if(i!=index ){
               autoObj.childNodes[i].className='suggest_link';
            }else{
               autoObj.childNodes[i].className='suggest_link_over';
               obj.value=autoObj.childNodes[i].innerHTML;
            }
        }
}