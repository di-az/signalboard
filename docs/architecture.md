# Overall Explanation
## How It Works
Signalboard is composed of a runtime that manages plugins, schedules their
execution, caches their output, and composes the resulting widgets into
documents that can be consumed by different display types.

### Configuration
The user defines which widgets should be displayed, their configuration,
and when they should be active. A display can have windows, which can have
multiple widgets. Widgets can be active during different time ranges.

### Engine
The engine coordinates the runtime. It initializes the configured
components and connects the scheduler, executor, and display composer.

### Scheduler
The scheduler determines which widgets should execute and when. It is
responsible for execution frequency and display cycles based on the user's
configuration.

### Plugins
Plugins provide the functionality and presentation for widgets. A plugin
can expose multiple widgets, each potentially providing a different way
of presenting the same underlying information.

Plugins communicate with the runtime through a JSON contract. Plugin
implementations should not be coupled to Signalboard's implementation
language.

### Executor
The executor handles plugin execution. When a widget needs to be
refreshed, it determines whether cached output can be reused or whether
the plugin needs to be executed.

The executor is responsible for invoking plugins, handling their results,
and updating the cache.

### Cache
Plugin output is cached by Signalboard. Cached output can be reused when
a widget is refreshed without requiring another plugin execution.

### Display Composer
The display composer takes the currently active widget outputs and
constructs the document that will be sent to a display.

Signalboard does not interpret the content of individual widgets. Widgets
control their own HTML and inline CSS, while the composer is responsible
for placing those widgets into the configured document.

### Displays
Displays consume the composed output.

A web display can consume the HTML document directly. An e-ink display
can render the document through Chromium and send the resulting frame to
the physical device.

The display itself does not need to know anything about plugins, widgets,
schedules, or their underlying data sources.


# Diagram
                       ┌──────────────────────┐                             
                       │    Configuration     │                             
                       └───────────┬──────────┘                             
                                   ▼                                        
     ┌──────────────────────────────────────────────────────────────┐       
     │                      Signalboard Runtime                     │       
     │                                                              │       
     │                                                              │       
     │                                                              │       
     │                        ┌─────────┐                           │       
     │               ┌────────┼ Engine  │          ┌────────┐       │       
     │               │        └────┬────┘          │ Plugin │       │       
     │               │             │               └───┬────┘       │       
     │               │             │                   │            │       
     │               │       ┌─────┴─────┐        ┌────┴─────┐      │       
     │               │       │ Scheduler ├────────┤ Executor │      │       
     │               │       └───────────┘        └────┬─────┘      │       
     │               │                                 │            │       
     │               │       ┌───────────┐         ┌───┴────┐       │       
     │               └──────►│  Display  │         │ Cache  │       │       
     │                       │  Composer │         └────────┘       │       
     │                       └─────┬─────┘                          │       
     │                             │                                │       
     │                             │                                │       
     └─────────────────────────────┼────────────────────────────────┘       
                                   │                                        
           HTML document           │                                        
                    ┌──────────────┴──────────────┐                         
                    ▼                             ▼                         
             ┌──────────────┐             ┌──────────────┐                  
             │ Web Display  │             │ E-ink Render │                  
             │              │             │              │                  
             │    HTML      │             │   Chromium   │                  
             └──────────────┘             │      ↓       │                  
                                          │    Frame     │                  
                                          └──────────────┘                  
