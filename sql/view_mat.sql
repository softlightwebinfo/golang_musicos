create materialized view if not exists view_mat_blogs as
SELECT b.id,
       b.title,
       b.description,
       b.fk_user_id,
       b.image,
       b.created_at,
       b.updated_at,
       b.slug,
       u.name       AS user_name,
       u.email      AS user_email,
       u.created_at AS user_created_at,
       u.updated_at AS user_updated_at,
       u.fk_role_id AS user_fk_role_id,
       u.phone      AS user_phone,
       u.active     AS user_active,
       u.age        AS user_age,
       bc.id        AS cat_id,
       bc.name      AS cat_name,
       bc.parent_id AS cat_parent_id
FROM ((blogs b
    JOIN users u ON ((b.fk_user_id = u.id)))
         JOIN blogs_categories bc ON ((b.fk_category_id = bc.id)))
WHERE (u.active = true);

alter materialized view view_mat_blogs owner to postgres;

create materialized view if not exists view_mat_categories_subcategories as
SELECT c.id                                                                 AS id_category,
       s.id                                                                 AS id_subcategory,
       c.name                                                               AS category,
       s.name                                                               AS subcategory,
       cs.slug,
       (SELECT count(i.id) AS count
        FROM items i
        WHERE ((i.fk_id_category = c.id) AND (i.fk_id_subcategory = s.id))) AS count_items
FROM ((categories_subcategories cs
    JOIN categories c ON ((cs.fk_id_category = c.id)))
         JOIN subcategories s ON ((cs.fk_id_subcategory = s.id)))
ORDER BY c.id, s.id;

alter materialized view view_mat_categories_subcategories owner to postgres;

create unique index if not exists view_mat_categories_subcategori_id_category_id_subcategory_idx1
    on view_mat_categories_subcategories (id_category, id_subcategory);

create index if not exists view_mat_categories_subcategorie_id_category_id_subcategory_idx
    on view_mat_categories_subcategories (id_category, id_subcategory);

create materialized view if not exists view_mat_items as
SELECT i.id              AS item_id,
       i.title           AS item_title,
       i.description     AS item_description,
       i.image           AS item_image,
       i.price           AS item_price,
       cs.category       AS item_category,
       cs.subcategory    AS item_subcategory,
       cs.id_category    AS item_id_category,
       cs.id_subcategory AS item_id_subcategory,
       i.contact_name    AS item_contact_name,
       i.contact_phone   AS item_contact_phone,
       i.created_at      AS item_created_at,
       i.updated_at      AS item_updated_at,
       u.id              AS user_id,
       u.name            AS user_name,
       u.email           AS user_email,
       u.age             AS user_age,
       u.created_at      AS user_created_at,
       u.updated_at      AS user_updated_at,
       u.active          AS user_active,
       u.phone           AS user_phone,
       cs.slug           AS item_category_slug,
       i.fk_city         AS item_city,
       c.country         AS item_country,
       c.region          AS item_region,
       c.name            AS city_name,
       r.name            AS region_name,
       co.name           AS country_name,
       co.slug           AS country_slug,
       r.slug            AS region_slug,
       c.slug            AS city_slug
FROM (((((items i
    JOIN view_mat_categories_subcategories cs ON (((i.fk_id_category = cs.id_category) AND
                                                   (i.fk_id_subcategory = cs.id_subcategory))))
    JOIN users u ON ((i.fk_user_id = u.id)))
    JOIN cities c ON ((i.fk_city = c.id)))
    JOIN regions r ON (((r.code = c.region) AND (r.country = c.country))))
         JOIN countries co ON ((r.country = co.code)))
WHERE (u.active = true)
ORDER BY i.updated_at DESC;

alter materialized view view_mat_items owner to postgres;

create materialized view if not exists view_mat_web as
SELECT t."primary",
       t.id,
       t.title,
       t.description,
       t.updated_at,
       t.type,
       t.slug,
       t.data
FROM (SELECT concat('blog-', b.id) AS "primary",
             b.id,
             b.title,
             b.description,
             b.updated_at,
             'blog'::text          AS type,
             b.slug                AS slug,
             NULL::text            AS data
      FROM view_mat_blogs b
      UNION
      SELECT concat('item-', i.item_id) AS "primary",
             i.item_id,
             i.item_title,
             i.item_description,
             i.item_updated_at,
             'item'::text               AS type,
             i.item_category_slug       AS slug,
             NULL::text                 AS data
      FROM view_mat_items i
      UNION
      SELECT concat('forum-', f.id) AS "primary",
             f.id,
             f.title,
             f.description,
             f.updated_at,
             'forum'::text          AS type,
             f.slug                 AS slug,
             NULL::text             AS data
      FROM forum f
      UNION
      SELECT concat('stores-', s.id) AS "primary",
             s.id,
             s.name,
             s.description,
             s.updated_at,
             'stores'::text          AS type,
             s.slug,
             s.web                   AS data
      FROM stores s) t
ORDER BY t.updated_at DESC;

alter materialized view view_mat_web owner to postgres;

