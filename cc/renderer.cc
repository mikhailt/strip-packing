#include "strip_packing.h"

static double max_x;
static double max_y;
static const double X = 1270;
static const double Y = 970;

double cx(double x) {
  return (x / max_x) * X;
}

double cy(double y) {
  return ((max_y - y) / max_y) * Y;
}

double ch(double h) {
  return (h / max_y) * Y;
}

double cw(double w) {
  return cx(w);
}

void DrawRect(cairo_t* cairo, Rect r) {
  cairo_rectangle(cairo, cx(r.x), cy(r.y + r.h), cw(r.w), ch(r.h));
}

gboolean expose_event(GtkWidget *widget, GdkEventExpose *event, gpointer data) {
  Context* context = (Context*)(data);
  cairo_t* cairo = gdk_cairo_create(widget->window);
  cairo_set_source_rgb(cairo, 0, 0, 0);
  cairo_set_line_width(cairo, 1);
  cairo_set_source_rgb(cairo, 0, 0, 0);

  max_x = std::max<double>(context->opt.m, 5);
  max_y = context->algo->solution_height;

  for (int y = 0; y < context->opt.m; ++y) {
    DrawRect(cairo, Rect(y, 0, 1, max_y));
  }
  cairo_stroke(cairo);
  std::vector<SavedRect>* vec = &context->algo->saved_rects;
  for (std::vector<SavedRect>::iterator i = vec->begin(); i != vec->end(); ++i) {
    if (i->color == "red") {
      cairo_set_source_rgb(cairo, 255, 0, 0);
    } else {
      cairo_set_source_rgb(cairo, 0, 0, 0);
    }
    DrawRect(cairo, i->r);
    cairo_stroke(cairo);
  }
 
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
