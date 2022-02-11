// Iframe Report
hlive.afterMessage.push(function() {
    document.querySelectorAll("[__iframeAttrReload__]").forEach(function (el) {
        const val = el.getAttribute("__iframeAttrReload__")

        if (val !== "<USED>") {
            el.setAttribute("__iframeAttrReload__", "<USED>")

            el.contentWindow.location.reload();
        }
    });
});
