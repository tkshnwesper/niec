$(document).ready(function() {
    hljs.initHighlightingOnLoad();
    $("#app").html(emojione.toImage($("#app").html()));
    $('[data-toggle="tooltip"]').tooltip();
    $('article table').each(function() {
        $(this).addClass("table table-bordered table-hover");
    });
    $('article img').each(function() {
        $(this).addClass("img-responsive center-block");
    })
});