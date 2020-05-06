<!DOCTYPE html>
<html>
<head>
    <title>Home</title>
</head>
<body>
<ul>
    {{ with .Players }}
        {{ range . }}
            <li>{{ .Name }}</li>
        {{ end }}
    {{ end }}
    <li></li>
</ul>
</body>
</html>
