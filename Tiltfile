version_settings(constraint='>=0.36.3')

def get_app_deps(name):
    return [
        "./apps/%s" % name,
        "./pnpm-workspace.yaml",
        "./pnpm-lock.yaml",
        "./package.json",
        "./packages",
    ]

def get_app_ignores(name):
    return [
        "./apps/%s/dist" % name,
        "./apps/%s/.turbo" % name,
    ]


# cluster - manifests
k8s_yaml([
    'infrastructure/k8s/app.namespace.yaml',
    'infrastructure/k8s/app.role.yaml',
])

# cluster - cluster setup
k8s_resource(
    objects=['app:namespace'],
    new_name='cluster-setup',
    labels=['cluster']
)

# cluster - app setup
k8s_resource(
    objects=['job-reader:role:app', 'job-reader-binding:rolebinding:app'],
    new_name='app-setup',
    labels=['cluster'],
    resource_deps=['cluster-setup']
)

# cluster - ghcr secret
local_resource(
    'ghcr-secret',
    cmd='''kubectl create secret docker-registry ghcr-secret \
        --dry-run=client \
        --docker-server=ghcr.io \
        --docker-username=$GITHUB_USERNAME \
        --docker-password=$GITHUB_TOKEN \
        --docker-email=$GITHUB_EMAIL \
        --namespace=app \
        -o yaml | kubectl apply -f -''',
    labels=['cluster'],
    resource_deps=['cluster-setup']
)

# cluster - packages setup
local_resource(
    'packages-setup',
    cmd='pnpm turbo build:packages',
    ignore=['packages/**/.turbo/**', 'packages/**/dist/**'],
    labels=['cluster']
)

# api - build
docker_build(
    ref="ghcr.io/awwwkshay/node-ts-api-starter", 
    context=".",
    dockerfile="./apps/api/Dockerfile",
    live_update=[
        sync("./apps/api/src/", "/app/apps/api/src"),
        run(
            "pnpm install",
            trigger=["./apps/api/package.json", "./apps/api/pnpm-lock.yaml"],
        )
    ],
    only=get_app_deps("api"),
    ignore=get_app_ignores("api"),
    target="development"
)

# api - manifests
k8s_yaml([
    'infrastructure/k8s/api/api.depl.yaml',
    'infrastructure/k8s/api/api.service.yaml',
    'infrastructure/k8s/api/api.config.yaml',
    'infrastructure/k8s/api/api.secret.yaml',
    'infrastructure/k8s/api/api.job.yaml',
    ])

# api - resource
k8s_resource(
    workload='api-deployment',
    objects=[
        'api-secret:secret:app',
        'api-config:configmap:app',
    ],
    port_forwards=['8000:8000','9000:9229'],
    resource_deps=['packages-setup', 'ghcr-secret'],
    labels=['api']
)

# api - job
k8s_resource(
    'api-job',
    labels=['api'],
    resource_deps=['packages-setup', 'ghcr-secret']
)

# web - build
docker_build(
    ref="ghcr.io/awwwkshay/node-ts-web-starter",
    context=".",
    dockerfile="./apps/web/Dockerfile",
    live_update=[
        sync("./apps/web/src/", "/app/apps/web/src"),
        sync("./apps/web/public/", "/app/apps/web/public"),
        run(
            "pnpm install",
            trigger=["./apps/web/package.json", "./apps/web/pnpm-lock.yaml"],
        )
    ],
    only=get_app_deps("web"),
    ignore=get_app_ignores("web"),
    target="development"
)

# web - manifests
k8s_yaml([  
    'infrastructure/k8s/web/web.depl.yaml',
    'infrastructure/k8s/web/web.service.yaml',
    'infrastructure/k8s/web/web.config.yaml',
    'infrastructure/k8s/web/web.secret.yaml',
    ])

# web - resource
k8s_resource(
    workload='web-deployment',
    objects=[
        'web-secret:secret:app',
        'web-config:configmap:app',
    ],
    port_forwards=['3000:3000','9001:9229'],
    resource_deps=['packages-setup', 'ghcr-secret'],
    labels=['web']
)

# auth - build
docker_build(
    ref="ghcr.io/awwwkshay/node-ts-auth-starter",
    context=".",
    dockerfile="./apps/auth/Dockerfile",
    live_update=[
        sync("./apps/auth/src/", "/app/apps/auth/src"),
        run(
            "pnpm install",
            trigger=["./apps/auth/package.json", "./apps/auth/pnpm-lock.yaml"],
        )
    ],
    only=get_app_deps("auth"),
    ignore=get_app_ignores("auth"),
    target="development"
)

# auth - manifests
k8s_yaml([
    'infrastructure/k8s/auth/auth.depl.yaml',
    'infrastructure/k8s/auth/auth.service.yaml',
    'infrastructure/k8s/auth/auth.config.yaml',
    'infrastructure/k8s/auth/auth.secret.yaml',
    'infrastructure/k8s/auth/auth.job.yaml', 
    ])

# auth - resource
k8s_resource(
    workload='auth-deployment',
    objects=[
        'auth-secret:secret:app',
        'auth-config:configmap:app',
    ],
    port_forwards=['4000:4000','9002:9229'],
    resource_deps=['packages-setup', 'ghcr-secret'],
    labels=['auth']
)

# auth - job
k8s_resource(
    'auth-job',
    labels=['auth'],
    resource_deps=['packages-setup', 'ghcr-secret']
)

# postgres - manifests
k8s_yaml([  
    'infrastructure/k8s/postgres/postgres.depl.yaml',
    'infrastructure/k8s/postgres/postgres.service.yaml',
    'infrastructure/k8s/postgres/postgres.secret.yaml',
    ])

# postgres - resource
k8s_resource(
    workload='postgres',
    objects=[
        'postgres-secret:secret:app',
    ],
    port_forwards=['5432:5432'],
    resource_deps=['ghcr-secret'],
    labels=['cluster']
)
