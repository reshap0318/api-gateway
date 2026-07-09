import type { DirectiveBinding } from 'vue'

interface ExtendedHTMLElement extends HTMLElement {
  _clickOutsideHandler?: (_event: MouseEvent) => void
}

export default {
  mounted(el: ExtendedHTMLElement, binding: DirectiveBinding<(_event: MouseEvent) => void>) {
    el._clickOutsideHandler = (event: MouseEvent) => {
      if (!(el === event.target || el.contains(event.target as Node))) {
        binding.value(event)
      }
    }
    document.addEventListener('click', el._clickOutsideHandler)
  },
  unmounted(el: ExtendedHTMLElement) {
    if (el._clickOutsideHandler) {
      document.removeEventListener('click', el._clickOutsideHandler)
      delete el._clickOutsideHandler
    }
  },
}
