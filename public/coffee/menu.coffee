set_active_tab = () ->
    for id in ["web", "cloud", "ip", "support"]
      req_str = "\/page\/" + id + "\/"
      regex = new RegExp(req_str)
      if (document.URL.match(regex))
        $("#" + id).addClass("active")
        break

$(document).ready ->
  set_active_tab()
