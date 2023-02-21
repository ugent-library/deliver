import htmx from 'htmx.org';
import BSN from 'bootstrap.native/dist/bootstrap-native-v4';

// Reinitialize Bootstrap Native after HTMX has udated the DOM.
export default function() {
  htmx.onLoad(function(el) {
    /*
     * ugly hack to prevent bootstrap from being initialized twice
     * context: htmx.onLoad is called once on page load, and then also
     * for every element fetched. But BSN.initCallback was already
     * called during page load, and this messes up the close button
     * of a toast.
     */
    if (el.tagName == 'BODY') return
    BSN.initCallback()
  });
}
