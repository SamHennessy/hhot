// Scroll To View
// Register plugin
if (hlive.afterMessage.get("hhui-scrollToView") === undefined) {
    hlive.afterMessage.set("hhui-scrollToView", function () {
        document.querySelectorAll("[hhui-scrollToView]").forEach(function (el) {
            el.scrollIntoView();
        });
    });
}