package pages

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/SamHennessy/hhot/hhotui/domain"
	"github.com/SamHennessy/hhot/hhotui/domain/topics"
	l "github.com/SamHennessy/hlive"
	"github.com/SamHennessy/hlive/hlivekit"
)

func Index(sl ServiceLocator) func() *l.Page {
	return func() *l.Page {
		p := l.NewPage()

		p.DOM().Head().Add(l.T("link",
			l.Attrs{"href": "https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css", "rel": "stylesheet"}))

		alloc := l.Box("")
		allocTotal := l.Box("")
		memRun := l.NewLockBox(true)
		allocH := l.CM("span", alloc, "/", allocTotal)
		allocH.SetMount(func(ctx context.Context) {
			go func() {
				for memRun.Get() {
					var mem runtime.MemStats
					runtime.ReadMemStats(&mem)
					alloc.Set(ByteSize(mem.Alloc))
					allocTotal.Set(ByteSize(mem.Sys))

					l.RenderComponent(ctx, allocH)

					time.Sleep(time.Second)
				}
			}()
		})
		allocH.SetUnmount(func(ctx context.Context) {
			memRun.Set(false)
		})

		logsDelete := l.C("button",
			l.Class("flex flex-shrink-0 p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.LogsClear, nil)
			}),
			iconBin, l.T("div", l.Class("hidden md:inline-block"), "Clear"),
		)

		logsMarker := l.C("button",
			l.Class("flex flex-shrink-0 p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.LogsAddMarker, nil)
			}),
			iconPageSep, l.T("div", l.Class("hidden md:inline-block"), "Spacer"),
		)

		restartBtn := l.C("button", l.Class("flex p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			iconRestart, l.T("div", l.Class("hidden md:inline-block"), "Restart"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.AppRestartDo, nil)
			}))

		stopBtn := l.C("button", l.Class("flex p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			iconStop, l.T("div", l.Class("hidden md:inline-block"), "Stop"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.AppStopDo, nil)
			}))

		startBtn := l.C("button", l.Class("flex p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			iconPlay, l.T("div", l.Class("hidden md:inline-block"), "Start"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.AppStartDo, nil)
			}))

		buildBtn := l.C("button", l.Class("flex p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			iconHammer, l.T("div", l.Class("hidden md:inline-block"), "Build"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.TriggerBuild, nil)
			}))

		refreshBtn := l.C("button", l.Class("flex p-1 text-gray-200 bg-gray-500 hover:bg-gray-700 hover:text-gray-300"),
			iconRefresh, l.T("div", l.Class("hidden md:inline-block"), "Refresh"),
			l.On("click", func(ctx context.Context, e l.Event) {
				sl.AppPubSub().Publish(topics.RefreshDo, nil)
			}))

		toolBox := l.C("div", l.Class("h-1/5 overflow-scroll overscroll-contain"),
			l.T("div",
				l.Class("absolute w-full flex bg-gradient-to-r from-gray-700 to-gray-500 items-center text-xs border-t-2 border-gray-600 h-7"),
				l.T("div", l.Class("w-12 text-gray-300 text-center"), "HHot "),
				l.T("div", l.Class("flex items-center flex-shrink-0"),
					l.T("div", l.Class("hidden md:inline-block text-gray-300 font-thin mx-1"), "Browser"),
					refreshBtn,
				),
				newBrowserPath(l.Class("flex-grow")),
				newBrowserTitle(),
			),
			l.T("div",
				l.Class("absolute w-full flex bg-gradient-to-r from-gray-700 to-gray-500 items-center text-xs h-6 mt-7 border-gray-600"),
				l.T("div", l.Class("w-20 text-gray-300 text-center font-thin"), allocH),
				l.T("div", l.Class("flex items-center flex-shrink-0"),
					l.T("div", l.Class("hidden md:inline-block text-gray-300 font-thin mx-1"), "App"),
					buildBtn, restartBtn, stopBtn, startBtn,
					l.T("div", l.Class("hidden md:inline-block text-gray-300 font-thin mx-1"), "Logs"),
					logsDelete, logsMarker,
				),
				l.T("div", l.Class("flex-grow")),
				l.T("div", buildLabel(sl)),
				l.T("div", appLabel(sl)),
			),
			l.T("div", l.Class("pt-14 pb-6 text-xs"),
				logList(sl),
			),
		)

		p.DOM().Body().Add(
			l.T("div", l.Class("h-screen"),
				l.T("div", l.Class("h-4/5"),
					buildMessage(sl),
					newAppIframe(sl),
				),
				toolBox,
			),
		)

		return p
	}
}

func buildLabel(sl ServiceLocator) l.Tagger {
	statusMessage := l.Box("")
	statusMessageC := l.T("div",
		l.Class("py-px px-1 ml-1 rounded-l md:w-20 flex"),
		iconHammer,
		l.T("div", l.Class("flex-grow")),
		l.T("div", l.Class("hidden md:inline-block"), statusMessage))

	statusTime := l.Box("")
	statusTimeC := l.T("div", l.Class("py-px px-0.5 mr-1 rounded-r bg-gray-700 text-gray-400 w-16 text-center"), statusTime)

	statusC := l.T("div", l.Class("flex flex-shrink-0 items-center"), statusMessageC, " ", statusTimeC)

	c := l.CM("div", l.Class("flex flex-shrink-0 items-center"), statusC)

	updateStatus := func(build domain.Build) {
		colourCSS := l.ClassBool{
			"bg-gray-400":   false,
			"bg-red-400":    false,
			"bg-yellow-400": false,
			"bg-green-400":  false,
			"text-gray-700": true,
		}

		switch build.Status {
		case domain.BuildStatusEmpty:
			statusMessage.Set("Waiting")
			colourCSS["bg-gray-400"] = true
		case domain.BuildStatusFailed:
			statusMessage.Set("Failed")
			colourCSS["bg-red-400"] = true
		case domain.BuildStatusInProgress:
			statusMessage.Set("Building")
			colourCSS["bg-yellow-400"] = true
		case domain.BuildStatusSuccess:
			statusMessage.Set("Success")
			colourCSS["bg-green-400"] = true
		}

		statusMessageC.Add(colourCSS)

		statusTime.Set(build.BuildDuration().Truncate(time.Millisecond).String())
	}

	// SSR
	updateStatus(sl.App().Build)

	var subFn hlivekit.QueueSubscriber

	c.SetMount(func(ctx context.Context) {
		onBuild := func(build domain.Build) {
			updateStatus(build)

			// Looking at sl closure
			go func() {
				for sl.App().Build.Status == domain.BuildStatusInProgress {
					updateStatus(sl.App().Build)
					l.Render(ctx)
					time.Sleep(time.Millisecond * 100)
				}
			}()
		}

		subFn = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			build, ok := message.Value.(domain.Build)
			if !ok {
				return
			}

			onBuild(build)
		})
		// Listen for updates
		sl.AppPubSub().Subscribe(subFn, topics.Build)

		// Init
		onBuild(sl.App().Build)
	})

	c.SetUnmount(func(ctx context.Context) {
		sl.AppPubSub().Unsubscribe(subFn, topics.Build)
	})

	return c
}

func appLabel(sl ServiceLocator) l.Tagger {
	statusMessage := l.Box("")
	statusMessageC := l.T("div",
		l.Class("py-px px-1 rounded-l md:w-20 flex"),
		iconPlay,
		l.T("div", l.Class("flex-grow")),
		l.T("div", l.Class("hidden md:inline-block"), statusMessage))

	statusTime := l.Box("")
	statusTimeC := l.T("div", l.Class("py-px px-0.5 mr-1 rounded-r bg-gray-700 text-gray-400 w-16 text-center"), statusTime)

	statusC := l.T("div", l.Class("flex flex-shrink-0 items-center"), statusMessageC, " ", statusTimeC)

	c := l.CM("div", l.Class("flex flex-shrink-0 items-center"), statusC)

	updateStatus := func(app *domain.App) {
		colourCSS := l.ClassBool{
			"bg-gray-400": false,
			"bg-red-400":  false,
			// TODO: use for starting
			"bg-yellow-400": false,
			"bg-green-400":  false,
			"text-gray-700": true,
		}

		switch app.Runtime.Status {
		case domain.RuntimeStatusEmpty:
			statusMessage.Set("Waiting")
			statusTime.Set("-s")
			colourCSS["bg-gray-400"] = true
		case domain.RuntimeStatusRunning:
			statusMessage.Set("Running")
			statusTime.Set(app.Runtime.Time().Truncate(time.Second).String())
			colourCSS["bg-green-400"] = true
		case domain.RuntimeStatusStopped:
			statusMessage.Set("Stopped")
			statusTime.Set("-s")
			colourCSS["bg-red-400"] = true
		}

		statusMessageC.Add(colourCSS)
	}

	// SSR
	updateStatus(sl.App())

	var subFn hlivekit.QueueSubscriber

	c.SetMount(func(ctx context.Context) {
		onBuild := func(app *domain.App) {
			updateStatus(app)

			go func() {
				for app.Runtime.Status == domain.RuntimeStatusRunning {
					updateStatus(sl.App())
					l.Render(ctx)
					time.Sleep(time.Second)
				}
			}()
		}

		subFn = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			app, ok := message.Value.(*domain.App)
			if !ok && app != nil {
				return
			}

			onBuild(app)
		})

		// Listen for updates
		sl.AppPubSub().Subscribe(subFn, topics.AppStart, topics.AppStop)

		// Init
		onBuild(sl.App())
	})

	c.SetUnmount(func(ctx context.Context) {
		sl.AppPubSub().Unsubscribe(subFn, topics.AppStart, topics.AppStop)
	})

	return c
}

func buildMessage(sl ServiceLocator) l.Componenter {
	c := l.CM("div")

	appBuild := sl.App().Build

	heading := l.Box("")
	output := l.Box("")

	c.Add(
		l.Class("flex flex-col"),
		l.T("div", l.Style{"height": "80vh"}, l.Class("overflow-scroll overscroll-contain"),
			l.T("div",
				l.Class("bg-gray-200 text-sm m-2 p-3 rounded-lg"),
				l.T("h1", l.Class("text-2xl m-5"), heading),
				l.T("pre", output),
			),
		),
	)

	update := func(build domain.Build) {
		output.Set(build.Output)

		switch build.Status {
		case domain.BuildStatusEmpty:
			heading.Set("Waiting for first build...")
		case domain.BuildStatusFailed:
			heading.Set("Build failed")
		case domain.BuildStatusSuccess:
			heading.Set("Build success")
		case domain.BuildStatusInProgress:
			heading.Set("Building...")
		}

		if build.Status == domain.BuildStatusFailed || build.Status == domain.BuildStatusInProgress {
			c.Add(l.ClassBool{"hidden": false})
		} else {
			c.Add(l.ClassBool{"hidden": true})
		}
	}

	// SSR
	update(appBuild)

	renderApp := func(app *domain.App) {
	}

	// SSR
	renderApp(sl.App())

	var (
		subFnBuild hlivekit.QueueSubscriber
		subFnApp   hlivekit.QueueSubscriber
	)

	c.SetMount(func(ctx context.Context) {
		subFnBuild = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			build, ok := message.Value.(domain.Build)
			if !ok {
				return
			}

			update(build)

			l.Render(ctx)
		})

		// Listen for updates
		sl.AppPubSub().Subscribe(subFnBuild, topics.Build)

		subFnApp = hlivekit.NewSub(func(message hlivekit.QueueMessage) {
			app, ok := message.Value.(*domain.App)
			if !ok {
				return
			}

			update(app.Build)

			renderApp(app)

			l.Render(ctx)
		})

		sl.AppPubSub().Subscribe(subFnApp, topics.AppStart, topics.AppStop)

		// Initial render
		update(sl.App().Build)
		renderApp(sl.App())
	})

	c.SetUnmount(func(ctx context.Context) {
		sl.AppPubSub().Unsubscribe(subFnBuild, topics.Build)
		sl.AppPubSub().Unsubscribe(subFnApp, topics.AppStart, topics.AppStop)
	})

	return c
}

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

// ByteSize returns a human-readable byte string of the form 10M, 12.5K, and so forth.  The following units are available:
//
//	E: Exabyte
//	P: Petabyte
//	T: Terabyte
//	G: Gigabyte
//	M: Megabyte
//	K: Kilobyte
//	B: Byte
//
// The unit that results in the smallest number greater than or equal to 1 is always chosen.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)

	switch {
	case bytes >= EXABYTE:
		unit = "E"
		value = value / EXABYTE
	case bytes >= PETABYTE:
		unit = "P"
		value = value / PETABYTE
	case bytes >= TERABYTE:
		unit = "T"
		value = value / TERABYTE
	case bytes >= GIGABYTE:
		unit = "G"
		value = value / GIGABYTE
	case bytes >= MEGABYTE:
		unit = "M"
		value = value / MEGABYTE
	case bytes >= KILOBYTE:
		unit = "K"
		value = value / KILOBYTE
	case bytes >= BYTE:
		unit = "B"
	case bytes == 0:
		return "0B"
	}

	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")

	return result + unit
}
