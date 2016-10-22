$(document).ready(function() {
    hljs.initHighlightingOnLoad();
    emojify.setConfig({
        emojify_tag_type : 'div',           // Only run emojify.js on this element
        only_crawl_id    : null,            // Use to restrict where emojify.js applies
        img_dir          : '/static/bower_components/emojify.js/dist/images/basic',  // Directory for emoji images
        ignored_tags     : {                // Ignore the following tags
            'SCRIPT'  : 1,
            'TEXTAREA': 1,
            'A'       : 1,
            'PRE'     : 1,
            'CODE'    : 1
        }
    });
    emojify.run();
    $('[data-toggle="tooltip"]').tooltip();
})