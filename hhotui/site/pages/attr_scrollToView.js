// Scroll To View
// Register plugin
hlive.afterMessage.push(function() {
    document.querySelectorAll("[hhui-scrollToView]").forEach(function (el) {
        el.scrollIntoView();
    });
});
