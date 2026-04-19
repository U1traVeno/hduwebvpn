// HTTP/1.1 200
// Content-Type: application/javascript
// Connection: keep-alive
// Date: Sun, 19 Apr 2026 08:47:33 GMT
// Last-Modified: Wed, 25 Jun 2025 09:25:46 GMT
// ETag: W/"685bc09a-1a96"
// Access-Control-Allow-Origin: *
// Access-Control-Allow-Methods: GET, POST, OPTIONS
// Access-Control-Allow-Headers: DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,user-header
// Strict-Transport-Security: max-age=15724800; includeSubDomains
// Server: nginx/1.251
// Content-Length: 6806

(function () {
    function loadStyleString(css) {
      var _style = document.createElement('style'),
          _head = document.head ? document.head : document.getElementsByTagName('head')[0];
      _style.type = 'text/css';
      try {
          _style.appendChild(document.createTextNode(css));
      } catch (ex) {
          _style.styleSheet.cssText = css;
      }
      _head.appendChild(_style);
      return _style;
    }

    loadStyleString('@media screen and (orientation: portrait) { .pl_box { width: 20px;height: 30px;background-color: red;opacity: 0;position: fixed;top: -20px;left: -20px;}} @media screen and (orientation: landscape) {.pl_box {width: 30px;height: 20px;background-color: green;opacity: 0;position: fixed;top: -20px;left: -20px;}}')

    var newElement = document.createElement('div');
    newElement.classList.add('pl_box');
    newElement.id ='pl_dom';
    document.body.appendChild(newElement);

    
    var neworientation = {}
    window.neworientation = neworientation;
    function debounce(func, delay) {
      let timer;
      return function () {
        clearTimeout(timer);
        timer = setTimeout(function() {
          func.apply(this, arguments);
        }, delay);
      };
    }
    var customEvent = document.createEvent('Event');
    customEvent.initEvent('pl_event', true, true);
    customEvent.detail = window.neworientation;

    const init = function(){
      const width = document.getElementById('pl_dom').clientWidth;
      const height = document.getElementById('pl_dom').clientHeight;
      if (parseInt(width) < parseInt(height)) {
        neworientation.current = 'portrait';
        neworientation.init = 'portrait';
      } else {
        neworientation.current = 'landscape';
        neworientation.init = 'landscape';
      }
      document.dispatchEvent(customEvent);
    }
    
    init();
    const handleResize = debounce(function() {
      init();
    }, 300);

    window.addEventListener('resize', handleResize, false);
  })();
var openAttestation = 'old'
var mediaQueryList = {};
mediaQueryList.matches = window.neworientation.init == 'portrait' ? true : false;
landscape = false;
portrait = false;
landscape = !mediaQueryList.matches;
portrait = mediaQueryList.matches;
var device = ''
outoinnerHeight = 0;
outoinnerWidth = 0;
function getDIr() {
    var baseUrl = document.getElementById('frontend-addr') && document.getElementById('frontend-addr').innerText;
    baseUrl = (baseUrl && baseUrl.replace(/^(https?:\/\/)?[^/]+/, '')) || '';
    try {
        window.frontend_base_url = baseUrl || staticBaseUrl;
    } catch (error) {
        window.frontend_base_url = '';
    }
    baseUrl = window.frontend_base_url;
    let rfurl = baseUrl + '/linkid/protected/api/dictconfig/get';
    var xhr = new XMLHttpRequest();
    // 请求接口的URL
    // 需要发送的数据
    const data = ['pc.login.page.config.type']
    xhr.open("POST", rfurl, true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.setRequestHeader('Csrf-Key', 'FzgxPikIetYDlXZM4lRG9taclVDa99lB');
    xhr.setRequestHeader('Csrf-Value', '7964f321f00366a3a287a133dd307ed0');
    // 发起POST请求
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // 请求成功
            var response = JSON.parse(xhr.responseText);
            if (response.data && response.data["pc.login.page.config.type"]) {
                openAttestation = response.data["pc.login.page.config.type"]
            }
            localStorage.setItem('openNewAttestation', response.data["pc.login.page.config.type"]);
            handleLandscape(response.data["pc.login.page.config.type"] || 'old')
        } else if (xhr.status !== 200 && xhr.readyState === 4 ) {
            // 请求失败
            console.error(xhr.statusText);
            localStorage.setItem('openNewAttestation', 'old');
            handleLandscape(false)
        }
    };
    xhr.send(JSON.stringify(data));
}
function handleLandscape(res){
    var url = document.getElementById("frontend-addr").innerText;
    var script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = url + '/public/deploy/deploy.js?' + (new Date()).getTime();
    if (landscape && res == 'new') {
        device = 'PC';
        script.onload = function () {
            window.casPageInit("loginnew");
        };

    } else if (landscape) {
        device = 'PC';
        script.onload = function () {
            window.casPageInit("login");
        };
    }
    else if (portrait) {
        device = 'PHONE';
        script.onload = function () {
            window.casPageInit("login");
        };
    }
    script.onerror = function (e) {
        console.log(e);
    };
    document.body.appendChild(script);
}
getDIr()
document.addEventListener('pl_event', function(event) {
    mediaQueryList.matches = window.neworientation.current == 'portrait' ? true : false;
    landscape = !mediaQueryList.matches;
    portrait = mediaQueryList.matches;

    if (portrait && !outoinnerHeight) {
        outoinnerHeight = window.innerHeight;
        outoinnerWidth = window.innerWidth;
    }
    if ((device !== 'PC' && outoinnerHeight && window.innerWidth === outoinnerWidth) || (window.innerWidth - window.innerHeight < 100 && window.innerWidth - window.innerHeight > 0))
    {
        device = 'PHONE';
    } else {
        if (device == 'PHONE' && !portrait && openAttestation == 'new') {
            location.reload();
        } else if (device == 'PC' && portrait && openAttestation == 'new') {
            location.reload();
        }
    }
});    


function viewport() {
    function isTabletLandscape() {
      // 获取屏幕的宽度和高度
      const screenWidth = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
      const screenHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
      // 判断是否为平板电脑
      const isTablet = screenWidth >= 768 && screenHeight <= 750; // 通常平板电脑的宽度大于等于768px
      // 判断是否为横屏模式
      const isLandscape = screenWidth > screenHeight;
      // 返回是否为平板电脑且横屏模式
      return isTablet && isLandscape;
    }

    function getAndZoomElement() {
      const element = document.getElementById('contentContainer');
      if (element) {
        element.style.zoom = '80%';
      } else {
        setTimeout(getAndZoomElement, 20); 
      }
    }

    //通过设备的屏幕宽度和高度，判断监听是pad设备，如果是则缩放比例设置为85%
    if (isTabletLandscape()) {
      if (window.navigator.userAgent.indexOf('HarmonyOS') > -1) {
        getAndZoomElement();
      }
    }
  };