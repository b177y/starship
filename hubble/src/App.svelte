<script>
    import { Router, Route } from "svelte-routing";
    import Notifications from 'svelte-notifications';
    import Home from "./pages/Home.svelte";
    import NetworkPage from "./pages/network.svelte";
    import NewNetworkPage from "./pages/newnet.svelte";
    import LoginPage from "./pages/login.svelte";
    import NetworkSidebar from "./sidebar.svelte";
    export let url = ""; //This property is necessary declare to avoid ignore the Router
    let networks = [];
    import API from './api.js';
    export const getAllNetworks = async() => {
        try {
            const response = await API.get("/networks/all");
            console.log(response)
            return response;
        } catch (error) {
            console.error(error);
        }
    };
    export const updateNetworks = async() => {
        console.log("Updating network list")
        networks = await getAllNetworks();
    }
    import { onDestroy } from 'svelte';
	import { LoggedIn } from './store.js';
  
	let showLogin;
	const unLoggedIn = LoggedIn.subscribe(v => showLogin = !v);
    console.log("showLogin", showLogin)
	onDestroy(unLoggedIn);
	
</script>

<main class="h-screen">
    <Notifications>
        {#if !showLogin}
        <div class="grid grid-cols-5 h-screen">
            <Router url="{url}">
                <NetworkSidebar networks={networks} on:updateNetworks={updateNetworks}/>
                <div class="col-span-4 bg-gray-100">
                    <Route path="/"><Home /></Route>
                    <Route path="networks/:netname" let:params>
                        <NetworkPage netname="{params.netname}" on:updateNetworks={updateNetworks}/>
                    </Route>
                    <Route path="networks/new">
                        <NewNetworkPage on:updateNetworks={updateNetworks}/>
                    </Route>
                </div>
            </Router>
        </div>
        {:else}
            <LoginPage />
        {/if}
    </Notifications>
</main>

<style global lang="postcss">
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>
