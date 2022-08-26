// Iframe Report
if (hlive.afterMessage.get("hhui-iframeReload") === undefined) {
    hlive.afterMessage.set("hhui-iframeReload", function () {
        document.querySelectorAll("[__iframeAttrReload__]").forEach(function (el) {
            const val = el.getAttribute("__iframeAttrReload__")

            if (val !== "<USED>") {
                el.setAttribute("__iframeAttrReload__", "<USED>")

                el.contentWindow.location.reload();
            }
        });
    });
}
