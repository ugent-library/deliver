export default function () {

  /*
   * Important: fallback document.execCommand("copy") does not work
   * here because we cannot call the method "select()" on the underlying "span"
   */

  let largeEvtHandler = function(evt) {

    let val = this.querySelector(".btn-clipboard-value").innerText
    let btn = this

    navigator.clipboard.writeText(val).then(function() {

      // clicking the button between green and original
      // makes it fail
      btn.removeEventListener("click", largeEvtHandler)

      let originalBtn = btn.cloneNode(true)
      btn.classList.remove("btn-outline-secondary")
      btn.classList.add("btn-outline-success")
      let icon = btn.querySelector("i")
      icon.removeAttribute("class")
      icon.classList.add("if", "if-check", "text-success")
      let span = btn.querySelector("span")
      span.removeAttribute("class")
      span.setAttribute("aria-live", "polite")
      span.innerText = "Copied link"

      setTimeout(function(){
        btn.parentNode.replaceChild(originalBtn, btn)
        originalBtn.addEventListener("click", largeEvtHandler)
      }, 1000)

    })
  }

  let smallEvtHandler = function(evt) {

    let val = this.querySelector(".btn-clipboard-value").innerText
    let btn = this

    navigator.clipboard.writeText(val).then(function() {

      // clicking the button between green and original
      // makes it fail
      btn.removeEventListener("click", smallEvtHandler)

      let originalBtn = btn.cloneNode(true)
      let icon = btn.querySelector("i")
      icon.removeAttribute("class")
      icon.classList.add("if", "if-check", "text-success")
      let span = btn.querySelector("span")
      span.classList.remove("text-primary")
      span.classList.add("text-success")
      span.setAttribute("aria-live", "polite")
      span.innerText = "Copied link"

      setTimeout(function(){
        btn.parentNode.replaceChild(originalBtn, btn)
        originalBtn.addEventListener("click", smallEvtHandler)
      }, 1000)

    })

  }

  document.querySelectorAll(".btn-clipboard-large").forEach(function(btn){
    btn.addEventListener("click", largeEvtHandler)
  })
  document.querySelectorAll(".btn-clipboard-small").forEach(function(btn){
    btn.addEventListener("click", smallEvtHandler)
  })
}
