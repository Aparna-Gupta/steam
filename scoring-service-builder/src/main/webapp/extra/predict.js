// predict.js
// H2O Scoring Service

// TODO
// check for timeouts in $.get calls

// Sort values if not sorted properly
// This is needed because c-1, c-2, ... are sorted alphabetically which is wrong
function sortValues(values) {
  if (values == null || values.length < 2) return values;

  var newvalues = [];
  if (values[0] == "c-1") {
    ok = true;
    // check if this is c-1, c-2 etc and make sure it's sorted correctly
    newvalues = [];
    for (i = 1; i <= values.length; i += 1) {
       if (values[i - 1] != "c-" + i) {
         ok = false;
       }
       newvalues.push("c-" + i);
    }
  }
  else {
    newvalues = values;
    ok = true;
  }
  result = !ok ? newvalues : newvalues.sort();
  return result;
}

// Form for model, shown in element
function showModel(model, element) {
    var names = model["m"]["_names"];
    var domains = model["m"]["_domains"];
    var domainMap = model["domainMap"];

    form = '<form id="allparams" action="" method="get">';
    for(i in names) {
        n = names[i];
        i1 = Number(i) + 1;
        form += '<p><b>' + i1 + '. ' + n + '</b> ';
        domain = domains[i];
        domain = sortValues(domain);
        card = domain == null ? 0 : domain.length;

        if (card < 2)
          form += '<input type="text" size=40 name="' + n + '">';
        else if (card <= 12) {
          for (i = 0; i < card; i += 1) {
            form += '<input type="radio" name="' + n + '" value="' + domain[i] + '"> ' + domain[i] + '</input>\n';
          }
        }
        else {
          form += '<select name="' + n + '">';
          form += '<option value=""></option>';
          for (v of domain) {
             form += '<option value="' + v + '">' + v + '</option>';
          }
          form += '</select>';
        }
        form += '</p>\n';
    }

    outputDomain = domains[i1];
//    console.log(outputDomain);

    form += '<input type="button" name="okbutton" value="PREDICT" onClick="runpred2(this.form)">';

    form += ' <input type="reset" value="CLEAR">';

    form += '</form>';
    element.innerHTML += form;

    form += '<p>';

}

function showInputParameters() {
    $.get('/info', function(data, status) {
    // show result
    info = document.querySelector(".input");
    showModel(data, info);
  },'json');
}

function showResult(div, status, data) {
    if ("classProbabilities" in data) {
        // binomial and multinomial
        label = data["label"];
        index = data["labelIndex"];
        probs = data["classProbabilities"];
        prob = probs[index];

        result = "Label <b>" + label + "</b> with probability <b>" + (prob * 100.0).toFixed(1) + "%</b>.<p>"
            + "Output Labels: [" + outputDomain + "]<br>"
            + "Class Probabilities: [" + probs + "]<br>"
            + "Label Index: " + index;
    }
    else if ("cluster" in data) {
        // clustering result
        result = "Cluster <b>" + data["cluster"] + "</b>";
    }
    else if ("value" in data) {
        // regression result
       result = "Value <b>" + data["value"] + "</b>";
    }
    else
        result = "Can't parse result";

    result += "<p><code>" + JSON.stringify(data) + "</code>";
    div.innerHTML = result;
    showStatistics();
}

function showUrl(pardiv, params) {
  // remove empty parameters returned by serialize.
  params = params.replace(/[\w]+=&/g, "").replace(/&?[\w]+=$/g, "");
  url = window.location.href + "predict?" + params;
  pardiv.innerHTML = '<a href="' + url + '" target="_blank"><code>' + url + '</code>';
}

function predResults(params) {
  pardiv = document.querySelector(".params");
  // add link that opens in new window
  showUrl(pardiv, params);

  div = document.querySelector(".results");
  cmd = '/predict?' + params;
    $.get(cmd, function(data, status) {
      showResult(div, status, data);
    },'json')
      .fail(function(data, status, error) {
//        console.log(arguments);
//        console.log("fail " + data.status + " " + data.statusText);
        down = "<b>Service is down</b>";
        div.innerHTML = down + "<br>status " + data.status + " statusText " + data.statusText;
        stats = document.querySelector(".stats");
        stats.innerHTML = down;
      });

}

function runpred(form) {
  predResults(form.p.value);
}

function runpred2(form) {
  predResults($('#allparams').serialize());
}

function showStats(div, data) {
    dayMs = 1000 * 60 * 60;
    upDays = Number(data['upTimeMs']) / dayMs;
    lastTimeAgoDays = Number(data['lastTimeAgoMs']) / dayMs;
    s = 'Service started ' + data['startTimeUTC'] + ', uptime ' + upDays.toFixed(3) + ' days. ';
    n = Number(data['numberOfPredictions']);
    if (n > 0) {
        s +=  '<br>'
        + 'Last prediction ' + data['lastTimeUTC'] + ', ' + lastTimeAgoDays.toFixed(3) + ' days ago.'
        +  '<br>'
        + 'Last prediction took ' + Number(data['lastPredictionMs']).toFixed(3) + ' ms. '
        +  '<br>'
        + data['numberOfPredictions'] + ' predictions in ' + Number(data['totalPredictionTimeMs']).toFixed(1)
        + ' (after ' + data['warmupFirstN'] + ' warmups ' + Number(data['totalPredictionTimeAfterWarmupMs']).toFixed(1) + ') ms. '
        +  '<br>'
        + 'Average prediction time ' + Number(data['avgPredictionTimeMs']).toFixed(3)
        + ' (after ' + data['warmupFirstN'] + ' warmups ' + Number(data['avgPredictionTimeAfterWarmupMs']).toFixed(3) + ') ms.'
        ;
    }
    url = window.location.href + "stats";
    s += '<p>More statistics at <code><a href="' + url + '" target="_blank">' + url + '</a>';
    div.innerHTML = s;
}

function showStatistics() {
    cmd = '/stats';
    res = $.get(cmd, function(data, status) {
        divs = document.querySelector(".stats");
        showStats(divs, data);
       },'json');
}


// main
showInputParameters();

showStatistics();




