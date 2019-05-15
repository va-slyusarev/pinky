let domain;
let uri;
let pinky;

let labelPinky = "Pinky";
let labelCopy = "Copy";

let received = false;

$(document).ready(function () {

    domain = $('#domain').val();
    uri = $('#uri');
    pinky = $('#pinky');

    uri.keypress(function (event) {
        let code = (event.keyCode ? event.keyCode : event.which);
        if (code == '13') {
            pinky.click()
        }
    });

    uri.bind('input', function () {
        if (!$.trim($(this).val())) {
            trigger(true)
        }
    });

    uri.focus();
    trigger(true);
});

function trigger(empty) {
    if (empty === true) {
        received = false;
        pinky.toggleClass("btn-outline-secondary", true).toggleClass("btn-outline-success", false);
        pinky.text(labelPinky);
    }

    if (received) {
        pinky.toggleClass("btn-outline-secondary", false).toggleClass("btn-outline-success", true);
        pinky.text(labelCopy);
    }
}

function Copy() {
    let value = $.trim(uri.val());

    if (!value) {
        return
    }

    uri.select();
    document.execCommand("copy");
    Success("The value copied to the clipboard: ", value);
}

function Pinky() {

    if (pinky.attr("disabled")) {
        return;
    }

    if (received) {
        Copy();
        return;
    }

    let value = $.trim(uri.val());

    if (!value) {
        Error("Please enter a non-empty uri value for Pinky.");
        return
    }

    pinky.attr('disabled', true);

    $.post({
        url: "/",
        data: value,
        success: function (msg, status) {
            if (msg === "") {
                Error("Internal Server Error");
                return;
            }
            received = true;
            Success("Success! Press Enter or Copy button to save value to the clipboard.");
            uri.val(domain + '/' + msg);
            pinky.attr('disabled', false);
            trigger()
        },
        error: function(xhr, textStatus, error){
            Error("Error: " + error);
        }
    })
}


function Success(message, uri) {
    showInfo(message, "alert-success", uri);
}

function Error(message, uri) {
    showInfo(message, "alert-secondary", uri);
}

let timeout;

function showInfo(message, type, uri) {
    clearTimeout(timeout);

    $('#info-placeholder').empty().append('<div id="info" class="alert ' + type + '">' + message + '<button class="close" data-dismiss="alert">&times;</button></div>')
    if (uri) {
        $('#info').append('<a href="' + uri + '" class="alert-link">' + uri + '</a>')
    }

    timeout = setTimeout(function () {
        $("#info").remove();
    }, 7000);
}