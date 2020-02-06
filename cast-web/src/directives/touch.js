import Hammer from 'hammerjs'

let ticking = false

let timer = null

let transform = {
  translate: { x: 0, y: 0 },
  scale: 1,
  angle: 0,
  rx: 0,
  ry: 0,
  rz: 0
}

function execute (f) {
  f && typeof f === 'function' && f()
}

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


function requestElementUpdate(el) {
  if(!ticking) {
      reqAnimationFrame(() => updateElementTransform(transform, el));
      ticking = true;
  }
}

function resetElement(el) {
  transform = {
      translate: { x: 0, y: 0 },
      scale: 1,
      angle: 0,
      rx: 0,
      ry: 0,
      rz: 0
  };
  requestElementUpdate(el);
}

const reqAnimationFrame = (function () {
  return window[Hammer.prefixed(window, 'requestAnimationFrame')] || function (callback) {
      window.setTimeout(callback, 1000 / 60);
  };
})();


function touch(el, binding = {}, vnode) {
  const defaultValues = { enable: false, limitClient: false, swipe: null, scale: false }
  const bindValues = binding.value || {}
  const options = { ...defaultValues, ...bindValues }
  console.log('options', options, bindValues, defaultValues)
  if (!options.enable) {
    el.style.transition = ''
    el.style.webkitTransform = '';
    el.style.mozTransform = '';
    el.style.transform = '';
    return
  }
  el.style.transition = 'transform 0.2s linear'
  let START_X = 0
  let START_Y = 0
  const CLIENT_W = window.innerWidth  || document.documentElement.clientWidth  || document.body.clientWidth;
  const CLIENT_H = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
  const maxTranslateX = CLIENT_W - el.offsetWidth - el.offsetLeft
  const maxTranslateY = CLIENT_H - el.offsetHeight - el.offsetTop
  const minTranslateX = el.offsetLeft * -1
  const minTranslateY = el.offsetTop * -1

  const mc = new Hammer.Manager(el)
  mc.add(new Hammer.Pan({ threshold: 0, pointers: 0 }))
  mc.add(new Hammer.Swipe()).recognizeWith(mc.get('pan'))
  mc.add(new Hammer.Tap({ event: 'doubletap', taps: 2 }))
  mc.add(new Hammer.Pinch({ threshold: 0 }))

  mc.on("hammer.input", function(ev) {
    if(ev.isFinal && transform.scale === 1) {
      START_X = transform.translate.x
      START_Y = transform.translate.y
    }
  })
  mc.on("swipe", () => {
    if (transform.scale === 1) {
      execute(binding.value && binding.value.swipe)
    }
  })
  mc.on("panmove", (ev) => {
    console.log('on move')
    let deltaX = START_X + ev.deltaX
    let deltaY = START_Y + ev.deltaY
    if (options.limitClient) {
      deltaX = Math.max(deltaX, minTranslateX)
      deltaX = Math.min(deltaX, maxTranslateX)
      deltaY = Math.max(deltaY, minTranslateY)
      deltaY = Math.min(deltaY, maxTranslateY)
    }
    transform.translate = {
      x: deltaX,
      y: deltaY
    };
    requestElementUpdate(el);
  });
  mc.on("doubletap", () => {
    if (!options.scale) { return }
    transform.scale = transform.scale === 1 ? 1.5 : 1
    transform.translate = {
      x: START_X,
      y: START_Y
    };
    requestElementUpdate(el)
  })
}


export default {
  install(Vue) {
    Vue.directive('touch', touch)
  }
}