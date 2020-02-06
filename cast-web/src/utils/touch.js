import Hammer from 'hammerjs'

function execute (f) {
  f && typeof f === 'function' && f()
}

const reqAnimationFrame = (function () {
  return window[Hammer.prefixed(window, 'requestAnimationFrame')] || function (callback) {
      window.setTimeout(callback, 1000 / 60);
  };
})();

const Touch = {
  bind(el, binding, vnode) {
    el.style.transition = 'transform 0.2s linear'
    const enable = binding.value.enable
    const START_X = el.offsetWidth
    const START_Y = el.offsetHeight
    let transform = {
      translate: { x: START_X, y: START_Y },
      scale: 1,
      angle: 0,
      rx: 0,
      ry: 0,
      rz: 0
    }
    let ticking = false
    let timer = null

    function updateElementTransform(transform, el) {
      let value = [
        'translate3d(' + transform.translate.x + 'px, ' + transform.translate.y + 'px, 0)',
        'scale(' + transform.scale + ', ' + transform.scale + ')',
        'rotate3d('+ transform.rx +','+ transform.ry +','+ transform.rz +','+  transform.angle + 'deg)'
      ];
    
      value = value.join(" ");
      el.style.webkitTransform = value;
      el.style.mozTransform = value;
      el.style.transform = value;
      ticking = false;
    }
    

    function requestElementUpdate() {
      if(!ticking) {
          reqAnimationFrame(() => updateElementTransform(transform, el));
          ticking = true;
      }
    }

    function resetElement() {
      transform = {
          translate: { x: START_X, y: START_Y },
          scale: 1,
          angle: 0,
          rx: 0,
          ry: 0,
          rz: 0
      };
      requestElementUpdate();
    }

    const mc = new Hammer.Manager(el)
    mc.add(new Hammer.Pan({ threshold: 0, pointers: 0 }));
    mc.add(new Hammer.Swipe()).recognizeWith(mc.get('pan'));
    mc.add(new Hammer.Tap({ event: 'doubletap', taps: 2 }))
    mc.add(new Hammer.Pinch({ threshold: 0 }))
    mc.on("hammer.input", function(ev) {
      if(ev.isFinal && transform.scale === 1) {
          resetElement();
      }
    })

    mc.on("swipe", () => {
      if (transform.scale === 1) {
        execute(binding.value.swipe)
      }
    })
    mc.on("panstart panmove", (ev) => {
      if (!enable) { return }
      transform.translate = {
        x: transform.translate.x + ev.deltaX,
        y: transform.translate.y + ev.deltaY
      };
      requestElementUpdate();
    });
    mc.on("doubletap", () => {
      if (!enable) { return }
      transform.scale = transform.scale === 1 ? 1.5 : 1
      transform.translate = {
        x: START_X,
        y: START_Y
      };
      requestElementUpdate()
    })
  }
}

export default {
  install(Vue) {
    Vue.directive('touch', Touch)
  }
}