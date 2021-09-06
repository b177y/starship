<script>
    import { link } from "svelte-routing";
    export let networks = [];
    import { onMount } from "svelte";
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    onMount(async () => {
        dispatch('updateNetworks', {node: node.node});
      });
    import { navigate } from "svelte-routing";
	import { LoggedIn, AuthToken } from './store.js';
    function logout(){
        AuthToken.update(_ => "")
        LoggedIn.update(_ => false)
        navigate('/');
    }
</script>

<main>
 <aside class="col-span-1 from-purple-600 bg-purple-800 h-screen text-white p-4 bg-gradient-to-tr top-0 sticky">
     <div class="flex flex-row justify-between align-end mb-4">
         <h1 class="text-xl">All Networks</h1>
         <button class="bg-purple-400 rounded-xl p-1"
                 on:click={logout}
         >Logout</button>
     </div>
        {#each networks as net}
            <a href="/networks/{net.name}" use:link class="no-underline hover:no-underline text-white hover:text-white visited:text-white link">
                <div class="shadow-xl hover:shadow-2xl my-2 bg-purple-900 p-2 rounded">
                    <h3 class="visited:text-white">Network: {net.name}</h3>
                    <h5>CIDR: {net.cidr}</h5>
                </div>
            </a>
        {/each}
        <a href="/networks/new" use:link class="w-full flex justify-center">
            <div class="shadow-xl hover:shadow-2xl my-2 bg-purple-500 p-2 rounded-xl w-6/12 flex justify-center">
                <h3 class="visited:text-white">Create New</h3>
            </div>
        </a>
 </aside>
</main>

<style>
    a {
        text-decoration: none;
    }
    a:visited {
      color: white;
    }
</style>
