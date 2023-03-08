export default function () {

  const restoreTimeout = 1500

  let evtHandler = function(evt){
    let value = this.dataset.value
    let btn = this

    navigator.clipboard.writeText(value).then(function() {
      btn.classList.remove("btn-outline-secondary")
      btn.classList.add("btn-outline-success")

      let icon = btn.querySelector("i")
      icon.classList.remove("if-copy", "text-primary")
      icon.classList.add("if-check", "text-success")

      let span = btn.querySelector("span")
      if (span) {
        span.classList.remove("text-primary")
        span.classList.add("text-success")
        span.setAttribute("aria-live", "polite")
        span.innerText="Link copied"
      }

      setTimeout(function(){
        icon.classList.remove("if-check", "text-success")
        icon.classList.add("if-copy", "text-primary")

        if (span) {
          span.classList.remove("text-success")
          span.classList.add("text-primary")
          span.innerText="Copy link"
        }
        
        btn.classList.remove("btn-outline-success")
        btn.classList.add("btn-outline-secondary")
      }, restoreTimeout)
    })
  }

  document.querySelectorAll(".btn-copy-to-clipboard").forEach(function(btn){
    btn.addEventListener("click", evtHandler)
  })

  document.querySelectorAll(".input-select-text").forEach(function(el){
    el.addEventListener("click", function(evt) {
      this.select()
    })
  })
}
