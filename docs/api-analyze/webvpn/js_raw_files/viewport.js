// HTTP/1.1 200
// Content-Type: application/javascript;charset=UTF-8
// Connection: keep-alive
// Date: Sun, 19 Apr 2026 08:47:33 GMT
// Strict-Transport-Security: max-age=15724800; includeSubDomains
// X-XSS-Protection: 1; mode=block
// Last-Modified: Fri, 27 Jun 2025 14:48:00 GMT
// Accept-Ranges: bytes
// Server: nginx/1.251
// Content-Length: 1250

(function () {
  var viewportMeta = document.querySelector('meta[name="viewport"]');
  if (!viewportMeta) {
    viewportMeta = document.createElement('meta');
    viewportMeta.name = 'viewport';
    document.head.appendChild(viewportMeta);
  }

  function isTabletLandscape() {
    // 获取屏幕的宽度和高度
    const screenWidth = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
    const screenHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
    const isTablet = screenWidth >= 768 && screenHeight <= 750; 
    // 判断是否为横屏模式
    const isLandscape = screenWidth > screenHeight;
    return isTablet && isLandscape;
  }

  //通过设备的屏幕宽度和高度，判断监听是pad设备，如果是则缩放比例设置为85%
  if(isTabletLandscape()){
    viewportMeta.content = 'width=device-width, initial-scale=0.8, maximum-scale=1, user-scalable=0';
    // if (window.navigator.userAgent.indexOf('HarmonyOS') > -1) {
    //   document.body.style.zoom = '80%';
    // } else {
    //   viewportMeta.content = 'width=device-width, initial-scale=0.8, maximum-scale=1, user-scalable=0';
    // }
  }
})()