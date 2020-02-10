import Hammer from 'hammerjs'

const ctx = '@@Touch';

let ticking = false

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

const defaultValues = { enable: false, limitClient: false, swipe: null, scale: false }

function touch(el, binding = {}, vnode) {
  const bindValues = binding.value || {}
  el[ctx] = { ...(el[ctx] || {}), options: { ...defaultValues, ...bindValues } }
  let START_X = 0
  let START_Y = 0
  let maxTranslateX
  let maxTranslateY
  let minTranslateX
  let minTranslateY

  const CLIENT_W = window.innerWidth  || document.documentElement.clientWidth  || document.body.clientWidth;
  const CLIENT_H = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
  

  vnode.context.$nextTick(() => {
    maxTranslateX = CLIENT_W - el.offsetWidth - el.offsetLeft
    maxTranslateY = CLIENT_H - el.offsetHeight - el.offsetTop
    minTranslateX = el.offsetLeft * -1
    minTranslateY = el.offsetTop * -1
  })

  if (!el[ctx].options.enable) {
    el.style.transition = ''
    el.style.webkitTransform = '';
    el.style.mozTransform = '';
    el.style.transform = '';
    if (el[ctx].mc) {
      el[ctx].mc.set({ enable: false });
    }
    return
  } else {
    el.style.transition = 'transform 0.2s linear'
  }

  if (el[ctx].mc) {
    el[ctx].mc.enable = true
    return
  }

  el[ctx].mc = new Hammer.Manager(el)
  el[ctx].mc.add(new Hammer.Pan({ threshold: 0, pointers: 0 }))
  el[ctx].mc.add(new Hammer.Swipe()).recognizeWith(el[ctx].mc.get('pan'))
  el[ctx].mc.add(new Hammer.Tap({ event: 'doubletap', taps: 2 }))
  el[ctx].mc.add(new Hammer.Pinch({ threshold: 0 }))


  el[ctx].mc.on("hammer.input", function(ev) {
    if(ev.isFinal && transform.scale === 1) {
      START_X = transform.translate.x
      START_Y = transform.translate.y
    }
  })

  el[ctx].mc.on("swipe", () => {
    if (transform.scale === 1) {
      execute(binding.value && binding.value.swipe)
    }
  })

  el[ctx].mc.on("panmove", (ev) => {
    let deltaX = START_X + ev.deltaX
    let deltaY = START_Y + ev.deltaY
    if (el[ctx].options.limitClient) {
      deltaX = minTranslateX ? Math.max(deltaX, minTranslateX) : deltaX
      deltaX = maxTranslateX ? Math.min(deltaX, maxTranslateX) : deltaX
      deltaY = minTranslateY ? Math.max(deltaY, minTranslateY) : deltaY
      deltaY = maxTranslateY ? Math.min(deltaY, maxTranslateY) : deltaY
    }
    transform.translate = {
      x: deltaX,
      y: deltaY
    };
    requestElementUpdate(el);
  });

  el[ctx].mc.on("doubletap", () => {
    if (!el[ctx].options.scale) { return }
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