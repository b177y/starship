<script>
    let nodes = []
    let network = {
        name: "",
        cidr: "",
        cipher: "",
        groups: [],
    }
    export let netname;
    import API from '../api.js';
    import { onMount } from 'svelte';
    import NodeCard from '../components/nodecard.svelte';
    import NetworkSettings from '../components/settings.svelte';
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher();
    export const getNetworkInfo = async() => {
        try {
            const response = await API.get(`/networks/${netname}/info`);
            console.log(response)
            return response;
        } catch (error) {
            console.error(error);
        }
    };
    export const getAllNodes = async() => {
        try {
            const response = await API.get(`/networks/${netname}/nodes/all`);
            console.log(response)
            return response;
        } catch (error) {
            console.error(error);
        }
    };

    async function updateNodes(){
        const res = await getAllNodes();
        nodes = res;
    }
    async function updateNetworkInfo(){
        const res = await getNetworkInfo();
        network = res;
    }
    async function refreshNet(){
        await updateNetworkInfo()
        await updateNodes()
    }
    onMount(async () => {
        console.log("Mounting component for network ", netname)
        const res = await getAllNodes();
        nodes = res;
        await updateNetworkInfo();
      });

    $: {
        console.log("Routing to ", netname);
        getAllNodes()
        .then(res => nodes = res)
        getNetworkInfo()
        .then(res => network = res)
    }

</script>

<main class="p-4 flex justify-center w-full">
    <div class="flex flex-col justify-center w-8/12">
        <div class="flex flex-row justify-between">
            <div class="flex flex-row items-end">
                <h1 class="text-2xl">{netname}</h1>
                <h2 class="text-gray-400 text-lg ml-4">{network.ca_fingerprint}</h2>
            </div>
            <button 
                on:click={async() => {await refreshNet()}}
                class="text-xl text-indigo-800 outline-none border-none">{'\u27F3'}</button>
        </div>
        <NetworkSettings netname={netname} network={network} on:updateNetworks={() => dispatch('updateNetworks')} on:updateNetworkInfo={async() => {await updateNetworkInfo()}}/>
        <h2 class="text-2xl mt-4">Nodes</h2>
        <div>
            {#each nodes as n}
                <NodeCard node={n} netname={netname} groups={network.groups} on:updateNodes={updateNodes}/>
            {/each}
        </div>
    </div>
</main>
