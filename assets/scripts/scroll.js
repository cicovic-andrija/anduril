$('#scrollToTop').removeClass('no-js');

$(window).scroll(function() {
    $(this).scrollTop() > 150
    ? $('#scrollToTop').fadeIn()
    : $('#scrollToTop').fadeOut();
});

$('#scrollToTop').click(function(e) {
    e.preventDefault();
    $("html, body").animate({scrollTop: 0}, "slow");
    return false;
});
