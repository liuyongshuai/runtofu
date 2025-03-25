{{ template "header" .}}
<div class="taglist">
    {{ range .tagList }}
        <span>
            <a href="/tag/{{ .TagId }}" target="_blank">{{ .TagName }}</a>
            <span class="badge">{{ .ContentNum }}</span>
        </span>
    {{ end }}
</div>
<script>
$( function(){for(var imageList="",i=0;34>=i;i++)var purl=STATIC_PREFIX+"/images/loading/"+i+".gif",imageList=imageList+('<img src="'+purl+'" class="hidden">');$("body").append(imageList);} );
</script>
{{ template "footer" .}}