<!DOCTYPE html>
<html>
<head>
    <title>Arctic Logs</title>
    <style>
        .column {
            float: left;
            width: 33.33%;
        }

        /* Clear floats after the columns */
        .row:after {
            content: "";
            display: table;
            clear: both;
        }

        .issues, .upload {
            font-size: 10pt;
        }

        .upload span {
            background: #eaeaea;
            color: #000000;
        }

        body {
            font-family: monospace;
            font-size: 12pt;
        }
    </style>
</head>
<body>
<p class="issues">See something that looks broken? Clan that shouldn't exist (etc.)? Go <a href="https://github.com/laser/arctic-logs-webservice/issues">here</a> and file an issue.</p>
<p class="upload">Have a log file you'd like added? Message <span>@matashiwa</span> on Discord.</p>
<div class="row">
    <div class="column">
        <u>Clans</u>
        <ul>
            {{ range .Clans}}
                <li><a href={{ .Url }}>{{ .Label }}</a></li>
            {{ end }}
        </ul>
    </div>
    <div class="column">
        <u>Players</u>
        <ul>
            {{ range .Players}}
                <li><a href={{ .Url }}>{{ .Label }}</a></li>
            {{ end }}
        </ul>
    </div>
    <div class="column">
        <u>Logs</u>
        <ul>
            {{ range .Logs}}
                <li><a href={{ .Url }}>{{ .Label }}</a></li>
            {{ end }}
        </ul>
    </div>
</div>

</body>
</html>


