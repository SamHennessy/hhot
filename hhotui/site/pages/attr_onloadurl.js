// Register plugin
if (hlive.beforeSendEvent.get("hhui-onloadurl") === undefined) {
    hlive.beforeSendEvent.set("hhui-onloadurl", function(event, message) {
        if (event.type !== "load") {
            return message
        }

        console.log(event.target.contentWindow.location, message);

        if (event.target.contentWindow && event.target.contentWindow.location) {
            if (!message.e) {
                message.e = {};
            }

            message.e.location = event.target.contentWindow.location.href;
        }

        return message
    });
}