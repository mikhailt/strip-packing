#include <string>

#include "strip_packing.h"

void configure_event(GtkWidget *widget, GdkEventExpose *event, gpointer data) {
  Context* context = (Context*)(data);
  if (NULL != context->ren->pixmap_) {
    g_object_unref(context->ren->pixmap_);
  }
  GtkAllocation allocation;
  gtk_widget_get_allocation(GTK_WIDGET(context->ren->drawing_area_), 
                            &allocation);
  GdkWindow* gdk_window = gtk_widget_get_window(
      GTK_WIDGET(context->ren->drawing_area_));
  context->ren->drawable_ = GDK_DRAWABLE(gdk_window);
  context->ren->pixmap_ = gdk_pixmap_new(context->ren->drawable_, 
                                         allocation.width, 
                                         allocation.height, 24);
  context->ren->gc_ = gdk_gc_new(GDK_DRAWABLE(context->ren->pixmap_));
  GdkColor color;
  gdk_color_parse("white", &color);
  gdk_gc_set_rgb_fg_color(context->ren->gc_, &color);
  gdk_draw_rectangle(context->ren->drawable_, context->ren->gc_, true, 0, 0, 
                     -1 , -1);
  gdk_gc_set_rgb_bg_color(context->ren->gc_, &color);
  gdk_color_parse("black", &color);
  gdk_gc_set_rgb_fg_color(context->ren->gc_, &color);
}

void Renderer::Init(Context* context) {
  pixmap_ = NULL;

  gtk_init(NULL, NULL);
  window_ = gtk_window_new(GTK_WINDOW_TOPLEVEL);
  g_signal_connect(GTK_OBJECT(window_), "destroy", gtk_main_quit, NULL);
  drawing_area_ = gtk_drawing_area_new();
  
  g_signal_connect(GTK_OBJECT(drawing_area_), "configure-event", 
                   G_CALLBACK(configure_event), context);
  
  
  vbox_ = gtk_vbox_new(true, 0);
  gtk_container_set_border_width(GTK_CONTAINER(vbox_), 5);
  
 /* 
  
  vbox := gtk.VBox(true, 0)
	vbox.SetBorderWidth(5)
	drawingarea = gtk.DrawingArea()

	var pixmap *gdk.GdkPixmap

	drawingarea.Connect("configure-event", func() {
		if pixmap != nil {
			pixmap.Unref()
		}
		var allocation gtk.GtkAllocation
		drawingarea.GetAllocation(&allocation)
		draw = drawingarea.GetWindow().GetDrawable()
		pixmap = gdk.Pixmap(drawingarea.GetWindow().GetDrawable(), allocation.Width, allocation.Height, 24)
		gc = gdk.GC(pixmap.GetDrawable())
		gc.SetRgbFgColor(gdk.Color("white"))
		pixmap.GetDrawable().DrawRectangle(gc, true, 0, 0, -1, -1)
		gc.SetRgbFgColor(gdk.Color("black"))
		gc.SetRgbBgColor(gdk.Color("white"))
	},
		nil)

	drawingarea.Connect("expose-event", func() {
		if pixmap != nil {
			drawingarea.GetWindow().GetDrawable().DrawDrawable(gc,
				pixmap.GetDrawable(), 0, 0, 0, 0, -1, -1)
		}
		draw_all(rects, m)
	},
		nil)

	vbox.Add(drawingarea)

	window.Add(vbox)
	window.Maximize()
	window.ShowAll()

	gtk.Main()*/
	
	gtk_container_add(GTK_CONTAINER(vbox_), drawing_area_);
	gtk_container_add(GTK_CONTAINER(window_), vbox_);
	gtk_window_maximize(GTK_WINDOW(window_));
	gtk_widget_show_all(GTK_WIDGET(window_));
	
	gtk_main();
}

void Renderer::DrawRects(std::vector<Rect>* rects, bool solid, std::string color) {
  
}