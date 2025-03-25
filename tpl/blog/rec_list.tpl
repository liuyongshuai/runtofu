{{ template "header" .}}
{{ .recListContent }}
<script>
$( function(){for(var imageList="",i=0;34>=i;i++)var purl=STATIC_PREFIX+"/images/loading/"+i+".gif",imageList=imageList+('<img src="'+purl+'" class="hidden">');$("body").append(imageList);} );
</script>
{{ template "footer" .}}