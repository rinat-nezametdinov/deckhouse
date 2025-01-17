<p align="center">
  <img src="https://github.com/deckhouse/deckhouse/blob/main/docs/site/images/d8-logo.png"/>
</p>

<p align="center">
  <a href="https://t.me/deckhouse"><img src="https://img.shields.io/badge/telegram-chat-179cde.svg?logo=telegram" alt="Telegram chat"></a>
  <a href="https://twitter.com/deckhouseio"><img src="https://img.shields.io/twitter/follow/deckhouseio?label=%40deckhouseio&style=flat-square" alt="Twitter"></a>
  <a href="https://github.com/deckhouse/deckhouse/discussions"><img src="https://img.shields.io/github/discussions/deckhouse/deckhouse" alt="GH Discussions"/></a>
</p>

[Deckhouse](https://deckhouse.io/) is an Open Source platform for managing Kubernetes clusters in a fully automatic and uniform fashion. It allows you to create homogeneous Kubernetes clusters anywhere and fully manages them. It supplies all the add-ons you need for auto-scaling, observability, security, and service mesh. It comes in Enterprise Edition (EE) and Community Edition (CE). [Certified in CNCF](https://github.com/cncf/k8s-conformance) for Kubernetes 1.19, 1.20, 1.21.

# Main features

- NoOps: system software on the nodes, Kubernetes core software, Kubernetes platform components are automatically managed.
- SLA by design: availability can be guaranteed even without direct access to your infrastructure.
- Completely identical and infrastructure-agnostic clusters. Deploy on a public cloud of your choice (AWS, GCP, Microsoft Azure, OVH Cloud), self-hosted cloud solutions (OpenStack and vSphere), and even bare-metal servers.
- 100 % vanilla Kubernetes based on an upstream version of Kubernetes.
- Easy to start: you need a couple of CLI commands and 8 minutes to get production-ready Kubernetes.
- A fully-featured platform. Many features *(check the diagram below)* — carefully configured & integrated — are available right out of the box.

A brief overview of essential Deckhouse Platform features, from infrastructure level to the platform:

<img src="https://github.com/deckhouse/deckhouse/blob/main/docs/site/images/diagrams/structure.svg?sanitize=true">

## CE vs. EE

While Deckhouse Platform CE is available free as an Open Source, EE is a commercial version of the platform that can be purchased with a paid subscription. EE's source is also open, but it's neither Open Source nor free to use.

EE brings many additional features that extend the basic functionality provided in CE. They include OpenStack & vSphere integration, Istio service mesh, multitenancy, enterprise-level security, BGP support, instant autoscaling, local DNS caching, and selectable timeframe for the platform's upgrades.

# Architecture

Deckhouse Platform follows the upstream version of Kubernetes, using that as a basis to build all of its features and configurations on. The added functionality is implemented via two building blocks:

- [shell-operator](https://github.com/flant/shell-operator) — to create Kubernetes operators *(please check the [KubeCon NA 2020 talk](https://www.youtube.com/watch?v=we0s4ETUBLc) for details)*;
- [addon-operator](https://github.com/flant/addon-operator) — to pack these operators into modules and manage them.

## Current status

Deckhouse Platform has a vast history of being used internally in Flant and is ready for production. Its beta testing started in May'21 when the first public demo tokens were issued. By the end of Jun'21, it was tested by 300+ engineers, and its source code went public via GitHub. The formal public announcement [was made](https://blog.flant.com/deckhouse-kubernetes-platform/) in the end of Jul'21.

Deckhouse Platform CE is now freely available for everyone. Deckhouse Platform EE can be accessed via 30-days tokens issued via [Deckhouse website](https://deckhouse.io/).

# Trying Deckhouse

Please, refer to the project's [Getting Started](https://deckhouse.io/en/gs/) to begin your journey with Deckhouse Platform. Choose the cloud provder or bare-metal option for your infrastructure and follow the relevant step-by-step instructions to deploy your first Deckhouse Kubernetes cluster.

If anything works in an unexpected manner or you have any questions, feel free to contact us via GitHub Issues / Discussions or reach a wider [community of Deckhouse users](#online-community) in Telegram and other resources.

# Online community

In addition to common GitHub features, here are some other online resources related to Deckhouse:

* [Twitter](https://twitter.com/deckhouseio) to stay informed about everything happening around Deckhouse;
* [Telegram chat](https://t.me/deckhouse) to discuss (there's a dedicated [Telegram chat in Russian](https://t.me/deckhouse_ru) as well);
* Flant's [tech blog](https://blog.flant.com/tag/deckhouse/) to read posts related to Deckhouse.
