function iosCopyToClipboard(el) {
  var oldContentEditable = el.contentEditable,
    oldReadOnly = el.readOnly,
    range = document.createRange();

  el.contentEditable = true;
  el.readOnly = false;
  range.selectNodeContents(el);

  var s = window.getSelection();
  s.removeAllRanges();
  s.addRange(range);

  el.setSelectionRange(0, 999999); // A big number, to cover anything that could be inside the element.

  el.contentEditable = oldContentEditable;
  el.readOnly = oldReadOnly;

  document.execCommand("copy");
}

function parseUrl(url) {
  var fakeLink = document.createElement("a");
  fakeLink.href = url;
  return fakeLink;
}

export default {
  copyEl: function(el) {
    let text = el.innerText;
    if (window) {
      if (navigator.userAgent.match(/ipad|ipod|iphone/i)) {
        iosCopyToClipboard(el);
      } else if (window.clipboardData && window.clipboardData.setData) {
        // IE specific code path to prevent textarea being shown while dialog is visible.
        try {
          window.clipboardData.setData("Text", text);
          return Promise.resolve();
        } catch (e) {
          return Promise.reject(e);
        }
      } else if (
        document.queryCommandSupported &&
        document.queryCommandSupported("copy")
      ) {
        var textarea = document.createElement("textarea");
        textarea.textContent = text;
        textarea.style.position = "fixed"; // Prevent scrolling to bottom of page in MS Edge.
        document.body.appendChild(textarea);
        textarea.select();
        try {
          document.execCommand("copy"); // Security exception may be thrown by some browsers.
          return Promise.resolve();
        } catch (ex) {
          console.warn("Copy to clipboard failed.", ex);
          return Promise.reject(ex);
        } finally {
          document.body.removeChild(textarea);
        }
      }
    }
  },
  reloadPage: function() {
    /* case 1: /#/aa/xx?abc=123           ---> /?t=TS#/aa/xx?abc=123
       case 2: /?abc=123#/aa/xx?abc=123   ---> /?t=TS&abc=123#/aa/xx?abc=123
       case 3: /aa/xx?abc=124             ---> /aa/xx?t=TS&abc=124
       case 4: /aa/xx                     ---> /aa/xx?t=TS
       case 5: /#/aa/xx                   ---> /?t=TS#/aa/xx
     */
    let url = window.location.href.replace(/token=[0-9a-zA-Z.-]+/, "");
    let tsQuery = `t=${Date.now()}`;
    let newUrl = url;
    if (url.indexOf("t=") !== -1) {
      // found timestamp
      newUrl = url.replace(/t=\d+/, tsQuery);
    } else {
      // not found
      let p = parseUrl(url);
      let query = "?" + tsQuery;
      if (p.search || p.search === "?") {
        query = p.search + "&" + tsQuery;
      }
      newUrl = `${p.protocol}//${p.hostname}${p.pathname}${query}${p.hash}`;
    }
    window.location.href = newUrl;
  }
};
