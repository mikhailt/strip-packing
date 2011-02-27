#include <string>

#include "strip_packing.h"

gboolean expose_event(GtkWidget *widget, GdkEventExpose *event, gpointer data) {
  Context* context = (Context*)(data);
  cairo_t* cairo = gdk_cairo_create(widget->window);
  cairo_set_source_rgb(cairo, 0, 0, 0);
  cairo_set_line_width (cairo, 0.1);
  cairo_move_to(cairo, 0, 0);
  cairo_line_to(cairo, 100, 100);
  cairo_stroke(cairo);
  cairo_destroy(cairo);
  gtk_widget_queue_draw(widget);
  return FALSE;
}

void Renderer::ShowAll(Context* context) {
  gtk_init(NULL, NULL);
  window_ = gtk_window_new(GTK_WINDOW_TOPLEVEL);
  g_signal_connect(window_, "expose-event", G_CALLBACK(expose_event), context);
  g_signal_connect(window_, "destroy", G_CALLBACK(gtk_main_quit), NULL);
  gtk_widget_set_app_paintable(window_, TRUE);
  gtk_window_maximize(GTK_WINDOW(window_));
  gtk_widget_show_all(window_);
  gtk_widget_queue_draw(window_);
  gtk_main();
}
