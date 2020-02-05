export default function downloadFile(url, fileName){
    var aLink = document.createElement('a');
    aLink.style.display = 'none';
    aLink.download = fileName;
    aLink.href = url
    document.body.appendChild(aLink);
    aLink.click()
    document.body.removeChild(aLink);
}