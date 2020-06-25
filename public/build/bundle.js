
(function(l, r) { if (l.getElementById('livereloadscript')) return; r = l.createElement('script'); r.async = 1; r.src = '//' + (window.location.host || 'localhost').split(':')[0] + ':35729/livereload.js?snipver=1'; r.id = 'livereloadscript'; l.getElementsByTagName('head')[0].appendChild(r) })(window.document);
var app = (function () {
    'use strict';

    function noop() { }
    function add_location(element, file, line, column, char) {
        element.__svelte_meta = {
            loc: { file, line, column, char }
        };
    }
    function run(fn) {
        return fn();
    }
    function blank_object() {
        return Object.create(null);
    }
    function run_all(fns) {
        fns.forEach(run);
    }
    function is_function(thing) {
        return typeof thing === 'function';
    }
    function safe_not_equal(a, b) {
        return a != a ? b == b : a !== b || ((a && typeof a === 'object') || typeof a === 'function');
    }

    function append(target, node) {
        target.appendChild(node);
    }
    function insert(target, node, anchor) {
        target.insertBefore(node, anchor || null);
    }
    function detach(node) {
        node.parentNode.removeChild(node);
    }
    function element(name) {
        return document.createElement(name);
    }
    function text(data) {
        return document.createTextNode(data);
    }
    function space() {
        return text(' ');
    }
    function listen(node, event, handler, options) {
        node.addEventListener(event, handler, options);
        return () => node.removeEventListener(event, handler, options);
    }
    function prevent_default(fn) {
        return function (event) {
            event.preventDefault();
            // @ts-ignore
            return fn.call(this, event);
        };
    }
    function attr(node, attribute, value) {
        if (value == null)
            node.removeAttribute(attribute);
        else if (node.getAttribute(attribute) !== value)
            node.setAttribute(attribute, value);
    }
    function children(element) {
        return Array.from(element.childNodes);
    }
    function set_style(node, key, value, important) {
        node.style.setProperty(key, value, important ? 'important' : '');
    }
    function custom_event(type, detail) {
        const e = document.createEvent('CustomEvent');
        e.initCustomEvent(type, false, false, detail);
        return e;
    }

    let current_component;
    function set_current_component(component) {
        current_component = component;
    }
    function get_current_component() {
        if (!current_component)
            throw new Error(`Function called outside component initialization`);
        return current_component;
    }
    function onMount(fn) {
        get_current_component().$$.on_mount.push(fn);
    }

    const dirty_components = [];
    const binding_callbacks = [];
    const render_callbacks = [];
    const flush_callbacks = [];
    const resolved_promise = Promise.resolve();
    let update_scheduled = false;
    function schedule_update() {
        if (!update_scheduled) {
            update_scheduled = true;
            resolved_promise.then(flush);
        }
    }
    function add_render_callback(fn) {
        render_callbacks.push(fn);
    }
    let flushing = false;
    const seen_callbacks = new Set();
    function flush() {
        if (flushing)
            return;
        flushing = true;
        do {
            // first, call beforeUpdate functions
            // and update components
            for (let i = 0; i < dirty_components.length; i += 1) {
                const component = dirty_components[i];
                set_current_component(component);
                update(component.$$);
            }
            dirty_components.length = 0;
            while (binding_callbacks.length)
                binding_callbacks.pop()();
            // then, once components are updated, call
            // afterUpdate functions. This may cause
            // subsequent updates...
            for (let i = 0; i < render_callbacks.length; i += 1) {
                const callback = render_callbacks[i];
                if (!seen_callbacks.has(callback)) {
                    // ...so guard against infinite loops
                    seen_callbacks.add(callback);
                    callback();
                }
            }
            render_callbacks.length = 0;
        } while (dirty_components.length);
        while (flush_callbacks.length) {
            flush_callbacks.pop()();
        }
        update_scheduled = false;
        flushing = false;
        seen_callbacks.clear();
    }
    function update($$) {
        if ($$.fragment !== null) {
            $$.update();
            run_all($$.before_update);
            const dirty = $$.dirty;
            $$.dirty = [-1];
            $$.fragment && $$.fragment.p($$.ctx, dirty);
            $$.after_update.forEach(add_render_callback);
        }
    }
    const outroing = new Set();
    let outros;
    function transition_in(block, local) {
        if (block && block.i) {
            outroing.delete(block);
            block.i(local);
        }
    }
    function transition_out(block, local, detach, callback) {
        if (block && block.o) {
            if (outroing.has(block))
                return;
            outroing.add(block);
            outros.c.push(() => {
                outroing.delete(block);
                if (callback) {
                    if (detach)
                        block.d(1);
                    callback();
                }
            });
            block.o(local);
        }
    }

    const globals = (typeof window !== 'undefined'
        ? window
        : typeof globalThis !== 'undefined'
            ? globalThis
            : global);
    function create_component(block) {
        block && block.c();
    }
    function mount_component(component, target, anchor) {
        const { fragment, on_mount, on_destroy, after_update } = component.$$;
        fragment && fragment.m(target, anchor);
        // onMount happens before the initial afterUpdate
        add_render_callback(() => {
            const new_on_destroy = on_mount.map(run).filter(is_function);
            if (on_destroy) {
                on_destroy.push(...new_on_destroy);
            }
            else {
                // Edge case - component was destroyed immediately,
                // most likely as a result of a binding initialising
                run_all(new_on_destroy);
            }
            component.$$.on_mount = [];
        });
        after_update.forEach(add_render_callback);
    }
    function destroy_component(component, detaching) {
        const $$ = component.$$;
        if ($$.fragment !== null) {
            run_all($$.on_destroy);
            $$.fragment && $$.fragment.d(detaching);
            // TODO null out other refs, including component.$$ (but need to
            // preserve final state?)
            $$.on_destroy = $$.fragment = null;
            $$.ctx = [];
        }
    }
    function make_dirty(component, i) {
        if (component.$$.dirty[0] === -1) {
            dirty_components.push(component);
            schedule_update();
            component.$$.dirty.fill(0);
        }
        component.$$.dirty[(i / 31) | 0] |= (1 << (i % 31));
    }
    function init(component, options, instance, create_fragment, not_equal, props, dirty = [-1]) {
        const parent_component = current_component;
        set_current_component(component);
        const prop_values = options.props || {};
        const $$ = component.$$ = {
            fragment: null,
            ctx: null,
            // state
            props,
            update: noop,
            not_equal,
            bound: blank_object(),
            // lifecycle
            on_mount: [],
            on_destroy: [],
            before_update: [],
            after_update: [],
            context: new Map(parent_component ? parent_component.$$.context : []),
            // everything else
            callbacks: blank_object(),
            dirty
        };
        let ready = false;
        $$.ctx = instance
            ? instance(component, prop_values, (i, ret, ...rest) => {
                const value = rest.length ? rest[0] : ret;
                if ($$.ctx && not_equal($$.ctx[i], $$.ctx[i] = value)) {
                    if ($$.bound[i])
                        $$.bound[i](value);
                    if (ready)
                        make_dirty(component, i);
                }
                return ret;
            })
            : [];
        $$.update();
        ready = true;
        run_all($$.before_update);
        // `false` as a special case of no DOM component
        $$.fragment = create_fragment ? create_fragment($$.ctx) : false;
        if (options.target) {
            if (options.hydrate) {
                const nodes = children(options.target);
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.l(nodes);
                nodes.forEach(detach);
            }
            else {
                // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
                $$.fragment && $$.fragment.c();
            }
            if (options.intro)
                transition_in(component.$$.fragment);
            mount_component(component, options.target, options.anchor);
            flush();
        }
        set_current_component(parent_component);
    }
    class SvelteComponent {
        $destroy() {
            destroy_component(this, 1);
            this.$destroy = noop;
        }
        $on(type, callback) {
            const callbacks = (this.$$.callbacks[type] || (this.$$.callbacks[type] = []));
            callbacks.push(callback);
            return () => {
                const index = callbacks.indexOf(callback);
                if (index !== -1)
                    callbacks.splice(index, 1);
            };
        }
        $set() {
            // overridden by instance, if it has props
        }
    }

    function dispatch_dev(type, detail) {
        document.dispatchEvent(custom_event(type, Object.assign({ version: '3.23.2' }, detail)));
    }
    function append_dev(target, node) {
        dispatch_dev("SvelteDOMInsert", { target, node });
        append(target, node);
    }
    function insert_dev(target, node, anchor) {
        dispatch_dev("SvelteDOMInsert", { target, node, anchor });
        insert(target, node, anchor);
    }
    function detach_dev(node) {
        dispatch_dev("SvelteDOMRemove", { node });
        detach(node);
    }
    function listen_dev(node, event, handler, options, has_prevent_default, has_stop_propagation) {
        const modifiers = options === true ? ["capture"] : options ? Array.from(Object.keys(options)) : [];
        if (has_prevent_default)
            modifiers.push('preventDefault');
        if (has_stop_propagation)
            modifiers.push('stopPropagation');
        dispatch_dev("SvelteDOMAddEventListener", { node, event, handler, modifiers });
        const dispose = listen(node, event, handler, options);
        return () => {
            dispatch_dev("SvelteDOMRemoveEventListener", { node, event, handler, modifiers });
            dispose();
        };
    }
    function attr_dev(node, attribute, value) {
        attr(node, attribute, value);
        if (value == null)
            dispatch_dev("SvelteDOMRemoveAttribute", { node, attribute });
        else
            dispatch_dev("SvelteDOMSetAttribute", { node, attribute, value });
    }
    function prop_dev(node, property, value) {
        node[property] = value;
        dispatch_dev("SvelteDOMSetProperty", { node, property, value });
    }
    function validate_slots(name, slot, keys) {
        for (const slot_key of Object.keys(slot)) {
            if (!~keys.indexOf(slot_key)) {
                console.warn(`<${name}> received an unexpected slot "${slot_key}".`);
            }
        }
    }
    class SvelteComponentDev extends SvelteComponent {
        constructor(options) {
            if (!options || (!options.target && !options.$$inline)) {
                throw new Error(`'target' is a required option`);
            }
            super();
        }
        $destroy() {
            super.$destroy();
            this.$destroy = () => {
                console.warn(`Component was already destroyed`); // eslint-disable-line no-console
            };
        }
        $capture_state() { }
        $inject_state() { }
    }

    /* src/CreateMap.svelte generated by Svelte v3.23.2 */

    const { console: console_1 } = globals;
    const file = "src/CreateMap.svelte";

    function create_fragment(ctx) {
    	let main;
    	let div53;
    	let br0;
    	let t0;
    	let h2;
    	let t2;
    	let br1;
    	let t3;
    	let form;
    	let div3;
    	let div2;
    	let div1;
    	let div0;
    	let t5;
    	let input0;
    	let t6;
    	let div7;
    	let div6;
    	let div5;
    	let div4;
    	let t8;
    	let input1;
    	let t9;
    	let div16;
    	let div11;
    	let div10;
    	let div9;
    	let div8;
    	let t11;
    	let input2;
    	let t12;
    	let div15;
    	let div14;
    	let div13;
    	let div12;
    	let t14;
    	let input3;
    	let t15;
    	let small0;
    	let t17;
    	let br2;
    	let t18;
    	let div52;
    	let div17;
    	let button0;
    	let t20;
    	let div51;
    	let div21;
    	let div20;
    	let div19;
    	let div18;
    	let t22;
    	let input4;
    	let t23;
    	let small1;
    	let t25;
    	let hr0;
    	let t26;
    	let div30;
    	let div25;
    	let div24;
    	let div23;
    	let div22;
    	let t28;
    	let input5;
    	let t29;
    	let div29;
    	let div28;
    	let div27;
    	let div26;
    	let t31;
    	let input6;
    	let t32;
    	let small2;
    	let t34;
    	let hr1;
    	let t35;
    	let div34;
    	let div33;
    	let div32;
    	let div31;
    	let t37;
    	let select0;
    	let option0;
    	let option1;
    	let option2;
    	let t41;
    	let small3;
    	let t43;
    	let hr2;
    	let t44;
    	let div38;
    	let div37;
    	let div36;
    	let div35;
    	let t46;
    	let select1;
    	let option3;
    	let option4;
    	let option5;
    	let t50;
    	let small4;
    	let t52;
    	let hr3;
    	let t53;
    	let div42;
    	let div41;
    	let div40;
    	let div39;
    	let t55;
    	let select2;
    	let option6;
    	let option7;
    	let t58;
    	let small5;
    	let t60;
    	let hr4;
    	let t61;
    	let div44;
    	let div43;
    	let input7;
    	let t62;
    	let label;
    	let t64;
    	let small6;
    	let t66;
    	let hr5;
    	let t67;
    	let div49;
    	let div47;
    	let div46;
    	let div45;
    	let t69;
    	let input8;
    	let t70;
    	let small7;
    	let t72;
    	let div48;
    	let p;
    	let t74;
    	let div50;
    	let t75;
    	let br3;
    	let t76;
    	let input9;
    	let t77;
    	let button1;
    	let t79;
    	let link;
    	let mounted;
    	let dispose;

    	const block = {
    		c: function create() {
    			main = element("main");
    			div53 = element("div");
    			br0 = element("br");
    			t0 = space();
    			h2 = element("h2");
    			h2.textContent = "Create a New Map";
    			t2 = space();
    			br1 = element("br");
    			t3 = space();
    			form = element("form");
    			div3 = element("div");
    			div2 = element("div");
    			div1 = element("div");
    			div0 = element("div");
    			div0.textContent = "Map Name";
    			t5 = space();
    			input0 = element("input");
    			t6 = space();
    			div7 = element("div");
    			div6 = element("div");
    			div5 = element("div");
    			div4 = element("div");
    			div4.textContent = "Number of Rounds";
    			t8 = space();
    			input1 = element("input");
    			t9 = space();
    			div16 = element("div");
    			div11 = element("div");
    			div10 = element("div");
    			div9 = element("div");
    			div8 = element("div");
    			div8.textContent = "Round Time, Minutes";
    			t11 = space();
    			input2 = element("input");
    			t12 = space();
    			div15 = element("div");
    			div14 = element("div");
    			div13 = element("div");
    			div12 = element("div");
    			div12.textContent = "Seconds";
    			t14 = space();
    			input3 = element("input");
    			t15 = space();
    			small0 = element("small");
    			small0.textContent = "Leave empty or zero for no time limit.";
    			t17 = space();
    			br2 = element("br");
    			t18 = space();
    			div52 = element("div");
    			div17 = element("div");
    			button0 = element("button");
    			button0.textContent = "Show advanced settings";
    			t20 = space();
    			div51 = element("div");
    			div21 = element("div");
    			div20 = element("div");
    			div19 = element("div");
    			div18 = element("div");
    			div18.textContent = "Grace Distance (m)";
    			t22 = space();
    			input4 = element("input");
    			t23 = space();
    			small1 = element("small");
    			small1.textContent = "Guesses within this distance (in meters) will be awarded full points.";
    			t25 = space();
    			hr0 = element("hr");
    			t26 = space();
    			div30 = element("div");
    			div25 = element("div");
    			div24 = element("div");
    			div23 = element("div");
    			div22 = element("div");
    			div22.textContent = "Population Density %, Minimum";
    			t28 = space();
    			input5 = element("input");
    			t29 = space();
    			div29 = element("div");
    			div28 = element("div");
    			div27 = element("div");
    			div26 = element("div");
    			div26.textContent = "Maximum";
    			t31 = space();
    			input6 = element("input");
    			t32 = space();
    			small2 = element("small");
    			small2.textContent = "0% is ocean. 10% is barren road. With 20%, you will find signs of civilization. Anything above 50% is already very populated.";
    			t34 = space();
    			hr1 = element("hr");
    			t35 = space();
    			div34 = element("div");
    			div33 = element("div");
    			div32 = element("div");
    			div31 = element("div");
    			div31.textContent = "Panorama connectedness";
    			t37 = space();
    			select0 = element("select");
    			option0 = element("option");
    			option0.textContent = "always";
    			option1 = element("option");
    			option1.textContent = "never";
    			option2 = element("option");
    			option2.textContent = "any";
    			t41 = space();
    			small3 = element("small");
    			small3.textContent = "If you want to be able to always walk somewhere or if you want single-image ones.";
    			t43 = space();
    			hr2 = element("hr");
    			t44 = space();
    			div38 = element("div");
    			div37 = element("div");
    			div36 = element("div");
    			div35 = element("div");
    			div35.textContent = "Copyright";
    			t46 = space();
    			select1 = element("select");
    			option3 = element("option");
    			option3.textContent = "any";
    			option4 = element("option");
    			option4.textContent = "Google only";
    			option5 = element("option");
    			option5.textContent = "third party only";
    			t50 = space();
    			small4 = element("small");
    			small4.textContent = "If you want to see only Google panos or also include third party panos.";
    			t52 = space();
    			hr3 = element("hr");
    			t53 = space();
    			div42 = element("div");
    			div41 = element("div");
    			div40 = element("div");
    			div39 = element("div");
    			div39.textContent = "Source";
    			t55 = space();
    			select2 = element("select");
    			option6 = element("option");
    			option6.textContent = "outdoors only";
    			option7 = element("option");
    			option7.textContent = "any";
    			t58 = space();
    			small5 = element("small");
    			small5.textContent = "If you want to exclude panoramas inside businesses.";
    			t60 = space();
    			hr4 = element("hr");
    			t61 = space();
    			div44 = element("div");
    			div43 = element("div");
    			input7 = element("input");
    			t62 = space();
    			label = element("label");
    			label.textContent = "Show labels on map";
    			t64 = space();
    			small6 = element("small");
    			small6.textContent = "Check this if the map should tell you how places are called.";
    			t66 = space();
    			hr5 = element("hr");
    			t67 = space();
    			div49 = element("div");
    			div47 = element("div");
    			div46 = element("div");
    			div45 = element("div");
    			div45.textContent = "Location string";
    			t69 = space();
    			input8 = element("input");
    			t70 = space();
    			small7 = element("small");
    			small7.textContent = "Constrain the game to a specified area - enter a country, state, city, neighborhood, lake, or any other bounded area.  Does not yet affect scoring.";
    			t72 = space();
    			div48 = element("div");
    			p = element("p");
    			p.textContent = "Sorry, that does not seem like a valid bounding box on OSM Nominatim.";
    			t74 = space();
    			div50 = element("div");
    			t75 = space();
    			br3 = element("br");
    			t76 = space();
    			input9 = element("input");
    			t77 = space();
    			button1 = element("button");
    			button1.textContent = "Create Map";
    			t79 = space();
    			link = element("link");
    			add_location(br0, file, 152, 4, 5645);
    			add_location(h2, file, 154, 4, 5655);
    			add_location(br1, file, 156, 4, 5686);
    			attr_dev(div0, "class", "input-group-text");
    			add_location(div0, file, 163, 20, 5903);
    			attr_dev(div1, "class", "input-group-prepend");
    			add_location(div1, file, 162, 16, 5849);
    			attr_dev(input0, "type", "text");
    			attr_dev(input0, "class", "form-control");
    			attr_dev(input0, "id", "Name");
    			add_location(input0, file, 165, 16, 5987);
    			attr_dev(div2, "class", "input-group");
    			add_location(div2, file, 161, 12, 5807);
    			attr_dev(div3, "class", "form-group");
    			add_location(div3, file, 160, 8, 5770);
    			attr_dev(div4, "class", "input-group-text");
    			add_location(div4, file, 172, 20, 6215);
    			attr_dev(div5, "class", "input-group-prepend");
    			add_location(div5, file, 171, 16, 6161);
    			attr_dev(input1, "type", "number");
    			attr_dev(input1, "class", "form-control");
    			attr_dev(input1, "id", "NumRounds");
    			input1.value = "5";
    			attr_dev(input1, "min", "1");
    			attr_dev(input1, "max", "100");
    			add_location(input1, file, 174, 16, 6307);
    			attr_dev(div6, "class", "input-group");
    			add_location(div6, file, 170, 12, 6119);
    			attr_dev(div7, "class", "form-group");
    			add_location(div7, file, 169, 8, 6082);
    			attr_dev(div8, "class", "input-group-text");
    			add_location(div8, file, 182, 24, 6610);
    			attr_dev(div9, "class", "input-group-prepend");
    			add_location(div9, file, 181, 20, 6552);
    			attr_dev(input2, "type", "number");
    			attr_dev(input2, "min", "0");
    			attr_dev(input2, "class", "form-control mr-sm-3");
    			attr_dev(input2, "id", "TimeLimit_minutes");
    			add_location(input2, file, 184, 20, 6713);
    			attr_dev(div10, "class", "input-group");
    			add_location(div10, file, 180, 16, 6506);
    			attr_dev(div11, "class", "col");
    			add_location(div11, file, 179, 12, 6472);
    			attr_dev(div12, "class", "input-group-text");
    			add_location(div12, file, 190, 24, 6988);
    			attr_dev(div13, "class", "input-group-prepend");
    			add_location(div13, file, 189, 20, 6930);
    			attr_dev(input3, "type", "number");
    			attr_dev(input3, "min", "0");
    			attr_dev(input3, "class", "form-control");
    			attr_dev(input3, "id", "TimeLimit_seconds");
    			add_location(input3, file, 192, 20, 7079);
    			attr_dev(div14, "class", "input-group");
    			add_location(div14, file, 188, 16, 6884);
    			attr_dev(div15, "class", "col");
    			add_location(div15, file, 187, 12, 6850);
    			attr_dev(div16, "class", "form-row");
    			add_location(div16, file, 178, 8, 6437);
    			attr_dev(small0, "class", "form-text text-muted");
    			add_location(small0, file, 196, 8, 7219);
    			add_location(br2, file, 200, 8, 7333);
    			attr_dev(button0, "class", "btn btn-info");
    			attr_dev(button0, "type", "button");
    			add_location(button0, file, 204, 16, 7433);
    			attr_dev(div17, "class", "card-header");
    			add_location(div17, file, 203, 12, 7391);
    			attr_dev(div18, "class", "input-group-text");
    			add_location(div18, file, 213, 28, 7933);
    			attr_dev(div19, "class", "input-group-prepend");
    			add_location(div19, file, 212, 24, 7871);
    			attr_dev(input4, "type", "number");
    			attr_dev(input4, "class", "form-control");
    			attr_dev(input4, "id", "GraceDistance");
    			input4.value = "10";
    			attr_dev(input4, "min", "0");
    			add_location(input4, file, 215, 24, 8043);
    			attr_dev(div20, "class", "input-group");
    			add_location(div20, file, 211, 20, 7821);
    			attr_dev(div21, "class", "form-group");
    			add_location(div21, file, 210, 16, 7776);
    			attr_dev(small1, "class", "form-text text-muted");
    			add_location(small1, file, 218, 16, 8191);
    			add_location(hr0, file, 221, 16, 8359);
    			attr_dev(div22, "class", "input-group-text");
    			add_location(div22, file, 227, 32, 8668);
    			attr_dev(div23, "class", "input-group-prepend");
    			add_location(div23, file, 226, 28, 8602);
    			attr_dev(input5, "type", "number");
    			attr_dev(input5, "class", "form-control mr-sm-3");
    			attr_dev(input5, "id", "MinDensity");
    			input5.value = "15";
    			attr_dev(input5, "min", "0");
    			attr_dev(input5, "max", "100");
    			add_location(input5, file, 229, 28, 8797);
    			attr_dev(div24, "class", "input-group");
    			add_location(div24, file, 225, 24, 8548);
    			attr_dev(div25, "class", "col");
    			add_location(div25, file, 224, 20, 8506);
    			attr_dev(div26, "class", "input-group-text");
    			add_location(div26, file, 235, 32, 9134);
    			attr_dev(div27, "class", "input-group-prepend");
    			add_location(div27, file, 234, 28, 9068);
    			attr_dev(input6, "type", "number");
    			attr_dev(input6, "class", "form-control mr-sm-3");
    			attr_dev(input6, "id", "MaxDensity");
    			input6.value = "100";
    			attr_dev(input6, "min", "0");
    			attr_dev(input6, "max", "100");
    			add_location(input6, file, 237, 28, 9241);
    			attr_dev(div28, "class", "input-group");
    			add_location(div28, file, 233, 24, 9014);
    			attr_dev(div29, "class", "col");
    			add_location(div29, file, 232, 20, 8972);
    			attr_dev(div30, "class", "form-row");
    			add_location(div30, file, 223, 16, 8463);
    			attr_dev(small2, "class", "form-text text-muted");
    			add_location(small2, file, 241, 16, 9436);
    			add_location(hr1, file, 245, 16, 9661);
    			attr_dev(div31, "class", "input-group-text");
    			add_location(div31, file, 250, 28, 9841);
    			attr_dev(div32, "class", "input-group-prepend");
    			add_location(div32, file, 249, 24, 9779);
    			option0.__value = "1";
    			option0.value = option0.__value;
    			option0.selected = "selected";
    			add_location(option0, file, 253, 28, 10032);
    			option1.__value = "2";
    			option1.value = option1.__value;
    			add_location(option1, file, 254, 28, 10112);
    			option2.__value = "0";
    			option2.value = option2.__value;
    			add_location(option2, file, 255, 28, 10172);
    			attr_dev(select0, "class", "form-control");
    			attr_dev(select0, "id", "Connectedness");
    			add_location(select0, file, 252, 24, 9955);
    			attr_dev(div33, "class", "input-group");
    			add_location(div33, file, 248, 20, 9729);
    			attr_dev(div34, "class", "form-group");
    			add_location(div34, file, 247, 16, 9684);
    			attr_dev(small3, "class", "form-text text-muted");
    			add_location(small3, file, 259, 16, 10302);
    			add_location(hr2, file, 263, 16, 10484);
    			attr_dev(div35, "class", "input-group-text");
    			add_location(div35, file, 268, 28, 10664);
    			attr_dev(div36, "class", "input-group-prepend");
    			add_location(div36, file, 267, 24, 10602);
    			option3.__value = "0";
    			option3.value = option3.__value;
    			option3.selected = "selected";
    			add_location(option3, file, 271, 28, 10838);
    			option4.__value = "1";
    			option4.value = option4.__value;
    			add_location(option4, file, 272, 28, 10915);
    			option5.__value = "2";
    			option5.value = option5.__value;
    			add_location(option5, file, 273, 28, 10980);
    			attr_dev(select1, "class", "form-control");
    			attr_dev(select1, "id", "Copyright");
    			add_location(select1, file, 270, 24, 10765);
    			attr_dev(div37, "class", "input-group");
    			add_location(div37, file, 266, 20, 10552);
    			attr_dev(div38, "class", "form-group");
    			add_location(div38, file, 265, 16, 10507);
    			attr_dev(small4, "class", "form-text text-muted");
    			add_location(small4, file, 277, 16, 11122);
    			add_location(hr3, file, 281, 16, 11293);
    			attr_dev(div39, "class", "input-group-text");
    			add_location(div39, file, 286, 28, 11473);
    			attr_dev(div40, "class", "input-group-prepend");
    			add_location(div40, file, 285, 24, 11411);
    			option6.__value = "1";
    			option6.value = option6.__value;
    			option6.selected = "selected";
    			add_location(option6, file, 289, 28, 11641);
    			option7.__value = "0";
    			option7.value = option7.__value;
    			add_location(option7, file, 290, 28, 11728);
    			attr_dev(select2, "class", "form-control");
    			attr_dev(select2, "id", "Source");
    			add_location(select2, file, 288, 24, 11571);
    			attr_dev(div41, "class", "input-group");
    			add_location(div41, file, 284, 20, 11361);
    			attr_dev(div42, "class", "form-group");
    			add_location(div42, file, 283, 16, 11316);
    			attr_dev(small5, "class", "form-text text-muted");
    			add_location(small5, file, 294, 16, 11858);
    			add_location(hr4, file, 298, 16, 12009);
    			attr_dev(input7, "class", "form-check-input");
    			attr_dev(input7, "type", "checkbox");
    			attr_dev(input7, "id", "ShowLabels");
    			input7.checked = true;
    			add_location(input7, file, 302, 24, 12144);
    			attr_dev(label, "class", "form-check-label");
    			attr_dev(label, "for", "label");
    			add_location(label, file, 303, 24, 12241);
    			attr_dev(div43, "class", "form-check form-check-inline");
    			add_location(div43, file, 301, 20, 12077);
    			attr_dev(div44, "class", "form-group");
    			add_location(div44, file, 300, 16, 12032);
    			attr_dev(small6, "class", "form-text text-muted");
    			add_location(small6, file, 306, 16, 12378);
    			add_location(hr5, file, 310, 16, 12538);
    			attr_dev(div45, "class", "input-group-text");
    			add_location(div45, file, 315, 28, 12734);
    			attr_dev(div46, "class", "input-group-prepend");
    			add_location(div46, file, 314, 24, 12672);
    			attr_dev(input8, "type", "text");
    			attr_dev(input8, "class", "form-control mr-sm-3");
    			attr_dev(input8, "id", "locString");
    			attr_dev(input8, "placeholder", "Location");
    			add_location(input8, file, 317, 24, 12842);
    			attr_dev(div47, "class", "input-group");
    			add_location(div47, file, 313, 20, 12622);
    			attr_dev(small7, "class", "form-text text-muted");
    			add_location(small7, file, 319, 20, 13006);
    			attr_dev(p, "class", "card-text");
    			add_location(p, file, 323, 24, 13358);
    			attr_dev(div48, "class", "card bg-danger text-white mt-1");
    			attr_dev(div48, "id", "error-dialog");
    			div48.hidden = true;
    			add_location(div48, file, 322, 20, 13264);
    			attr_dev(div49, "class", "form-group");
    			add_location(div49, file, 312, 16, 12577);
    			attr_dev(div50, "id", "bounds-map");
    			set_style(div50, "width", "80%");
    			set_style(div50, "height", "50vh");
    			set_style(div50, "margin-left", "10%");
    			set_style(div50, "margin-right", "10%");
    			add_location(div50, file, 326, 16, 13519);
    			attr_dev(div51, "class", "card-body");
    			attr_dev(div51, "id", "advanced-settings");
    			div51.hidden = /*advancedHidden*/ ctx[1];
    			add_location(div51, file, 209, 12, 7689);
    			attr_dev(div52, "class", "card border-info");
    			add_location(div52, file, 202, 8, 7348);
    			add_location(br3, file, 330, 8, 13661);
    			attr_dev(input9, "id", "hidden-input");
    			attr_dev(input9, "type", "hidden");
    			attr_dev(input9, "name", "result");
    			input9.value = "";
    			add_location(input9, file, 332, 8, 13676);
    			attr_dev(button1, "id", "submit-button");
    			attr_dev(button1, "type", "submit");
    			attr_dev(button1, "class", "btn btn-primary");
    			set_style(button1, "margin-bottom", "2em");
    			add_location(button1, file, 334, 8, 13749);
    			attr_dev(form, "method", "post");
    			add_location(form, file, 158, 4, 5696);
    			attr_dev(link, "rel", "stylesheet");
    			attr_dev(link, "href", "static/leaflet/leaflet.css");
    			add_location(link, file, 337, 4, 13879);
    			attr_dev(div53, "class", "container");
    			add_location(div53, file, 150, 4, 5616);
    			add_location(main, file, 149, 0, 5605);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, main, anchor);
    			append_dev(main, div53);
    			append_dev(div53, br0);
    			append_dev(div53, t0);
    			append_dev(div53, h2);
    			append_dev(div53, t2);
    			append_dev(div53, br1);
    			append_dev(div53, t3);
    			append_dev(div53, form);
    			append_dev(form, div3);
    			append_dev(div3, div2);
    			append_dev(div2, div1);
    			append_dev(div1, div0);
    			append_dev(div2, t5);
    			append_dev(div2, input0);
    			append_dev(form, t6);
    			append_dev(form, div7);
    			append_dev(div7, div6);
    			append_dev(div6, div5);
    			append_dev(div5, div4);
    			append_dev(div6, t8);
    			append_dev(div6, input1);
    			append_dev(form, t9);
    			append_dev(form, div16);
    			append_dev(div16, div11);
    			append_dev(div11, div10);
    			append_dev(div10, div9);
    			append_dev(div9, div8);
    			append_dev(div10, t11);
    			append_dev(div10, input2);
    			append_dev(div16, t12);
    			append_dev(div16, div15);
    			append_dev(div15, div14);
    			append_dev(div14, div13);
    			append_dev(div13, div12);
    			append_dev(div14, t14);
    			append_dev(div14, input3);
    			append_dev(form, t15);
    			append_dev(form, small0);
    			append_dev(form, t17);
    			append_dev(form, br2);
    			append_dev(form, t18);
    			append_dev(form, div52);
    			append_dev(div52, div17);
    			append_dev(div17, button0);
    			append_dev(div52, t20);
    			append_dev(div52, div51);
    			append_dev(div51, div21);
    			append_dev(div21, div20);
    			append_dev(div20, div19);
    			append_dev(div19, div18);
    			append_dev(div20, t22);
    			append_dev(div20, input4);
    			append_dev(div51, t23);
    			append_dev(div51, small1);
    			append_dev(div51, t25);
    			append_dev(div51, hr0);
    			append_dev(div51, t26);
    			append_dev(div51, div30);
    			append_dev(div30, div25);
    			append_dev(div25, div24);
    			append_dev(div24, div23);
    			append_dev(div23, div22);
    			append_dev(div24, t28);
    			append_dev(div24, input5);
    			append_dev(div30, t29);
    			append_dev(div30, div29);
    			append_dev(div29, div28);
    			append_dev(div28, div27);
    			append_dev(div27, div26);
    			append_dev(div28, t31);
    			append_dev(div28, input6);
    			append_dev(div51, t32);
    			append_dev(div51, small2);
    			append_dev(div51, t34);
    			append_dev(div51, hr1);
    			append_dev(div51, t35);
    			append_dev(div51, div34);
    			append_dev(div34, div33);
    			append_dev(div33, div32);
    			append_dev(div32, div31);
    			append_dev(div33, t37);
    			append_dev(div33, select0);
    			append_dev(select0, option0);
    			append_dev(select0, option1);
    			append_dev(select0, option2);
    			append_dev(div51, t41);
    			append_dev(div51, small3);
    			append_dev(div51, t43);
    			append_dev(div51, hr2);
    			append_dev(div51, t44);
    			append_dev(div51, div38);
    			append_dev(div38, div37);
    			append_dev(div37, div36);
    			append_dev(div36, div35);
    			append_dev(div37, t46);
    			append_dev(div37, select1);
    			append_dev(select1, option3);
    			append_dev(select1, option4);
    			append_dev(select1, option5);
    			append_dev(div51, t50);
    			append_dev(div51, small4);
    			append_dev(div51, t52);
    			append_dev(div51, hr3);
    			append_dev(div51, t53);
    			append_dev(div51, div42);
    			append_dev(div42, div41);
    			append_dev(div41, div40);
    			append_dev(div40, div39);
    			append_dev(div41, t55);
    			append_dev(div41, select2);
    			append_dev(select2, option6);
    			append_dev(select2, option7);
    			append_dev(div51, t58);
    			append_dev(div51, small5);
    			append_dev(div51, t60);
    			append_dev(div51, hr4);
    			append_dev(div51, t61);
    			append_dev(div51, div44);
    			append_dev(div44, div43);
    			append_dev(div43, input7);
    			append_dev(div43, t62);
    			append_dev(div43, label);
    			append_dev(div51, t64);
    			append_dev(div51, small6);
    			append_dev(div51, t66);
    			append_dev(div51, hr5);
    			append_dev(div51, t67);
    			append_dev(div51, div49);
    			append_dev(div49, div47);
    			append_dev(div47, div46);
    			append_dev(div46, div45);
    			append_dev(div47, t69);
    			append_dev(div47, input8);
    			append_dev(div49, t70);
    			append_dev(div49, small7);
    			append_dev(div49, t72);
    			append_dev(div49, div48);
    			append_dev(div48, p);
    			append_dev(div51, t74);
    			append_dev(div51, div50);
    			append_dev(form, t75);
    			append_dev(form, br3);
    			append_dev(form, t76);
    			append_dev(form, input9);
    			append_dev(form, t77);
    			append_dev(form, button1);
    			append_dev(div53, t79);
    			append_dev(div53, link);

    			if (!mounted) {
    				dispose = [
    					listen_dev(button0, "click", /*click_handler*/ ctx[4], false, false, false),
    					listen_dev(input8, "change", /*locStringUpdated*/ ctx[3], false, false, false),
    					listen_dev(form, "submit", prevent_default(/*handleFormSubmit*/ ctx[2]), false, true, false)
    				];

    				mounted = true;
    			}
    		},
    		p: function update(ctx, [dirty]) {
    			if (dirty & /*advancedHidden*/ 2) {
    				prop_dev(div51, "hidden", /*advancedHidden*/ ctx[1]);
    			}
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(main);
    			mounted = false;
    			run_all(dispose);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function intById(id, fallback = 0) {
    	let input = document.getElementById(id);

    	if (input) {
    		if (!input.value) {
    			return fallback;
    		}

    		return parseInt(input.value, 10);
    	} else {
    		console.log("Couldn't find input '" + id + "', using fallback.");
    		return fallback;
    	}
    }

    function strById(id, fallback = "") {
    	let input = document.getElementById(id);

    	if (input) {
    		return input.value;
    	} else {
    		console.log("Couldn't find input '" + id + "', using fallback.");
    		return fallback;
    	}
    }

    // given Nominatim results, takes the most significant one with a polygon or
    // multipolygon and returns it as a turf.multiPolygon
    function geojsonFromNominatim(data) {
    	console.log("getting geojson...");

    	for (let i = 0; i < data.length; i++) {
    		let type = data[i].geojson.type.toLowerCase();

    		if (type === "multipolygon") {
    			return turf.multiPolygon(data[i].geojson.coordinates);
    		} else if (type === "polygon") {
    			return turf.multiPolygon([data[i].geojson.coordinates]);
    		}
    	}

    	console.log("No matching polygon!");
    	return null;
    }

    function instance($$self, $$props, $$invalidate) {
    	const NOMINATIM_URL = locStringEncoded => `https://nominatim.openstreetmap.org/search?q=${locStringEncoded}&polygon_geojson=1&limit=5&polygon_threshold=0.005&format=json`;

    	let mapSettings = {
    		Name: "",
    		Polygon: null,
    		Area: 0,
    		NumRounds: 0,
    		TimeLimit: 0,
    		GraceDistance: 10,
    		MinDensity: 0,
    		MaxDensity: 100,
    		Connectedness: 0,
    		Copyright: 0,
    		Source: 0,
    		ShowLabels: true
    	};

    	let locString = "";
    	let previewMap;
    	let previewPolyGroup;
    	let advancedHidden = true;

    	onMount(async () => {
    		$$invalidate(0, previewMap = L.map("bounds-map", { center: [0, 0], zoom: 1 }));

    		L.tileLayer("https://api.mapbox.com/styles/v1/jwlarocque/ckb5ngk3e1rq11lr5ztqekj5b/tiles/256/{z}/{x}/{y}?access_token=pk.eyJ1IjoiandsYXJvY3F1ZSIsImEiOiJja2I1bmRtMW0xNm16MnlxcDl3ZTh3cTBoIn0.CJk6LRRfs-Chwmm--NiBfw", {
    			attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OSM</a> contributors, <a href=\"https://wikitech.wikimedia.org/wiki/Wikitech:Cloud_Services_Terms_of_use\">Wikimedia Cloud Services</a>"
    		}).addTo(previewMap);

    		previewPolyGroup = L.layerGroup().addTo(previewMap);
    	});

    	// collates createmap form data into a JSON object, 
    	// then sends a newmap request to the server
    	function handleFormSubmit() {
    		mapSettings.Name = strById("Name");
    		mapSettings.NumRounds = intById("NumRounds");
    		mapSettings.GraceDistance = intById("GraceDistance");
    		mapSettings.MinDensity = intById("MinDensity");
    		mapSettings.MaxDensity = intById("MaxDensity");
    		mapSettings.Connectedness = intById("Connectedness");
    		mapSettings.Copyright = intById("Copyright");
    		mapSettings.Source = intById("Source");
    		let showLabelsInput = document.getElementById("ShowLabels");

    		if (showLabelsInput) {
    			mapSettings.ShowLabels = showLabelsInput.checked;
    		}

    		// read total TimeLimit
    		mapSettings.TimeLimit = 0;

    		mapSettings.TimeLimit += 60 * intById("TimeLimit_minutes");
    		mapSettings.TimeLimit += intById("TimeLimit_seconds");

    		// sanity check density fields
    		// TODO: nicer error messages than alerts
    		// TODO: check that population density in Polygon overlaps with the
    		//       specified range (otherwise we'll never be able to find good
    		//       panos.)
    		if (mapSettings.MinDensity > mapSettings.MaxDensity) {
    			alert("Max density must be greater than min density.");
    			return;
    		}

    		// TODO: evaluate challenge generation (to make sure mapSettings aren't so
    		//       specific that it takes a huge number of API requests to find good
    		//       panos)
    		// send to server /newmap
    		fetch("/newmap", {
    			method: "POST",
    			headers: { "Content-Type": "application/json" },
    			body: JSON.stringify(mapSettings)
    		}).then(console.log("mapSettings sent to server"));
    	}

    	function locStringUpdated() {
    		let old = locString;
    		let locStringInput = document.getElementById("locString");

    		if (locStringInput) {
    			locString = document.getElementById("locString").value;
    		}

    		if (old !== locString) {
    			updatePolygonFromLocString();
    		}
    	}

    	function showPolygonOnMap() {
    		previewPolyGroup.clearLayers();
    		let map_poly = L.geoJSON(mapSettings.Polygon).addTo(previewPolyGroup);
    		previewMap.fitBounds(map_poly.getBounds());
    	}

    	function updatePolygonFromLocString() {
    		if (locString === "" || !locString) {
    			mapSettings.Polygon = null;
    			return;
    		}

    		fetch(NOMINATIM_URL(encodeURI(locString.replace(" ", "+")))).then(response => response.json()).then(data => {
    			mapSettings.Polygon = geojsonFromNominatim(data);
    			mapSettings.Area = turf.area(mapSettings.Polygon);
    			showPolygonOnMap();
    		});
    	}

    	const writable_props = [];

    	Object.keys($$props).forEach(key => {
    		if (!~writable_props.indexOf(key) && key.slice(0, 2) !== "$$") console_1.warn(`<CreateMap> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;
    	validate_slots("CreateMap", $$slots, []);

    	const click_handler = () => {
    		$$invalidate(1, advancedHidden = !advancedHidden);

    		setTimeout(
    			function () {
    				previewMap.invalidateSize();
    			},
    			400
    		);
    	};

    	$$self.$capture_state = () => ({
    		onMount,
    		NOMINATIM_URL,
    		mapSettings,
    		locString,
    		previewMap,
    		previewPolyGroup,
    		advancedHidden,
    		handleFormSubmit,
    		intById,
    		strById,
    		locStringUpdated,
    		showPolygonOnMap,
    		updatePolygonFromLocString,
    		geojsonFromNominatim
    	});

    	$$self.$inject_state = $$props => {
    		if ("mapSettings" in $$props) mapSettings = $$props.mapSettings;
    		if ("locString" in $$props) locString = $$props.locString;
    		if ("previewMap" in $$props) $$invalidate(0, previewMap = $$props.previewMap);
    		if ("previewPolyGroup" in $$props) previewPolyGroup = $$props.previewPolyGroup;
    		if ("advancedHidden" in $$props) $$invalidate(1, advancedHidden = $$props.advancedHidden);
    	};

    	if ($$props && "$$inject" in $$props) {
    		$$self.$inject_state($$props.$$inject);
    	}

    	$$self.$$.update = () => {
    		if ($$self.$$.dirty & /*previewMap*/ 1) {
    			 window.globalMap = previewMap;
    		}
    	};

    	return [previewMap, advancedHidden, handleFormSubmit, locStringUpdated, click_handler];
    }

    class CreateMap extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance, create_fragment, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "CreateMap",
    			options,
    			id: create_fragment.name
    		});
    	}
    }

    /* src/App.svelte generated by Svelte v3.23.2 */

    const { console: console_1$1 } = globals;
    const file$1 = "src/App.svelte";

    // (26:1) {:else}
    function create_else_block(ctx) {
    	let h3;

    	const block = {
    		c: function create() {
    			h3 = element("h3");
    			h3.textContent = "404.  That's an error.";
    			add_location(h3, file$1, 26, 2, 762);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, h3, anchor);
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(h3);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_else_block.name,
    		type: "else",
    		source: "(26:1) {:else}",
    		ctx
    	});

    	return block;
    }

    // (24:59) 
    function create_if_block_2(ctx) {
    	let p;

    	const block = {
    		c: function create() {
    			p = element("p");
    			p.textContent = "Create a new challenge.";
    			add_location(p, file$1, 24, 2, 720);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, p, anchor);
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(p);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_2.name,
    		type: "if",
    		source: "(24:59) ",
    		ctx
    	});

    	return block;
    }

    // (22:53) 
    function create_if_block_1(ctx) {
    	let createmap;
    	let current;
    	createmap = new CreateMap({ $$inline: true });

    	const block = {
    		c: function create() {
    			create_component(createmap.$$.fragment);
    		},
    		m: function mount(target, anchor) {
    			mount_component(createmap, target, anchor);
    			current = true;
    		},
    		i: function intro(local) {
    			if (current) return;
    			transition_in(createmap.$$.fragment, local);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(createmap.$$.fragment, local);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			destroy_component(createmap, detaching);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block_1.name,
    		type: "if",
    		source: "(22:53) ",
    		ctx
    	});

    	return block;
    }

    // (20:1) {#if window.location.pathname === "/"}
    function create_if_block(ctx) {
    	let p;

    	const block = {
    		c: function create() {
    			p = element("p");
    			p.textContent = "Landing page.";
    			add_location(p, file$1, 20, 2, 568);
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, p, anchor);
    		},
    		i: noop,
    		o: noop,
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(p);
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_if_block.name,
    		type: "if",
    		source: "(20:1) {#if window.location.pathname === \\\"/\\\"}",
    		ctx
    	});

    	return block;
    }

    function create_fragment$1(ctx) {
    	let main;
    	let nav;
    	let span;
    	let t1;
    	let ul;
    	let div;
    	let li0;
    	let a0;
    	let t3;
    	let li1;
    	let a1;
    	let t5;
    	let current_block_type_index;
    	let if_block;
    	let current;
    	const if_block_creators = [create_if_block, create_if_block_1, create_if_block_2, create_else_block];
    	const if_blocks = [];

    	function select_block_type(ctx, dirty) {
    		if (window.location.pathname === "/") return 0;
    		if (window.location.pathname === "/createmap") return 1;
    		if (window.location.pathname === "/createchallenge") return 2;
    		return 3;
    	}

    	current_block_type_index = select_block_type();
    	if_block = if_blocks[current_block_type_index] = if_block_creators[current_block_type_index](ctx);

    	const block = {
    		c: function create() {
    			main = element("main");
    			nav = element("nav");
    			span = element("span");
    			span.textContent = "Earthwalker";
    			t1 = space();
    			ul = element("ul");
    			div = element("div");
    			li0 = element("li");
    			a0 = element("a");
    			a0.textContent = "Home";
    			t3 = space();
    			li1 = element("li");
    			a1 = element("a");
    			a1.textContent = "Source code";
    			t5 = space();
    			if_block.c();
    			attr_dev(span, "class", "navbar-brand");
    			add_location(span, file$1, 7, 2, 176);
    			attr_dev(a0, "class", "nav-link");
    			attr_dev(a0, "href", "/");
    			add_location(a0, file$1, 11, 5, 328);
    			attr_dev(li0, "class", "nav-item active");
    			add_location(li0, file$1, 10, 4, 294);
    			attr_dev(a1, "class", "nav-link");
    			attr_dev(a1, "href", "https://gitlab.com/glatteis/earthwalker");
    			add_location(a1, file$1, 14, 5, 407);
    			attr_dev(li1, "class", "nav-item");
    			add_location(li1, file$1, 13, 4, 380);
    			attr_dev(div, "class", "collapse navbar-collapse");
    			add_location(div, file$1, 9, 3, 251);
    			attr_dev(ul, "class", "navbar-nav");
    			add_location(ul, file$1, 8, 2, 224);
    			attr_dev(nav, "class", "navbar navbar-expand-sm navbar-light bg-light");
    			add_location(nav, file$1, 6, 1, 114);
    			attr_dev(main, "class", "svelte-329bu3");
    			add_location(main, file$1, 5, 0, 106);
    		},
    		l: function claim(nodes) {
    			throw new Error("options.hydrate only works if the component was compiled with the `hydratable: true` option");
    		},
    		m: function mount(target, anchor) {
    			insert_dev(target, main, anchor);
    			append_dev(main, nav);
    			append_dev(nav, span);
    			append_dev(nav, t1);
    			append_dev(nav, ul);
    			append_dev(ul, div);
    			append_dev(div, li0);
    			append_dev(li0, a0);
    			append_dev(div, t3);
    			append_dev(div, li1);
    			append_dev(li1, a1);
    			append_dev(main, t5);
    			if_blocks[current_block_type_index].m(main, null);
    			current = true;
    		},
    		p: noop,
    		i: function intro(local) {
    			if (current) return;
    			transition_in(if_block);
    			current = true;
    		},
    		o: function outro(local) {
    			transition_out(if_block);
    			current = false;
    		},
    		d: function destroy(detaching) {
    			if (detaching) detach_dev(main);
    			if_blocks[current_block_type_index].d();
    		}
    	};

    	dispatch_dev("SvelteRegisterBlock", {
    		block,
    		id: create_fragment$1.name,
    		type: "component",
    		source: "",
    		ctx
    	});

    	return block;
    }

    function instance$1($$self, $$props, $$invalidate) {
    	const writable_props = [];

    	Object.keys($$props).forEach(key => {
    		if (!~writable_props.indexOf(key) && key.slice(0, 2) !== "$$") console_1$1.warn(`<App> was created with unknown prop '${key}'`);
    	});

    	let { $$slots = {}, $$scope } = $$props;
    	validate_slots("App", $$slots, []);
    	$$self.$capture_state = () => ({ CreateMap });
    	 console.log(window.location.pathname);
    	return [];
    }

    class App extends SvelteComponentDev {
    	constructor(options) {
    		super(options);
    		init(this, options, instance$1, create_fragment$1, safe_not_equal, {});

    		dispatch_dev("SvelteRegisterComponent", {
    			component: this,
    			tagName: "App",
    			options,
    			id: create_fragment$1.name
    		});
    	}
    }

    const app = new App({
    	target: document.body
    });

    return app;

}());
//# sourceMappingURL=bundle.js.map
