set_active_tab = () ->
    for id in ["web", "cloud", "voip", "linux"]
      req_str = "\/page\/" + id + "\/"
      regex = new RegExp(req_str)
      if (document.URL.match(regex))
        $("#" + id).addClass("active")
        break

set_active_icon = () ->
    labels = 'web': 'Сложные веб приложения просто', 'cloud': 'Просто создайте свое облако', 'voip': 'IP телефония', 'linux': 'GNU/Linux аутсорс'
    for link, label of labels
      console.log(link + label)
      req_str = "\/page\/" + link + "\/"
      regex = new RegExp(req_str)
      if (document.URL.match(regex))
        $('#icon').attr('src', '/public/images/' + link + '_icon.png')
        $('#icon_text').html(label)
        break

$(document).ready ->
  set_active_tab()
  #set_active_icon()
